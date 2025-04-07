require "result"
require "openssl"
require "random/secure"

module Audiofeeler
  class DeployInventory
    def initialize(@db : DB::Database)
      @db = db
    end

    def find_one(id)
      @db.query_one "SELECT id, account_id, server, local_dir, remote_dir, username, password FROM deploys WHERE id = ?", id do |rs|
        return Ok.done(
          Deploy.new(
            id: rs.read(Int64),
            account_id: rs.read(Int64),
            server: rs.read(String),
            local_dir: rs.read(String?),
            remote_dir: rs.read(String?),
            username: rs.read(String?),
            password: rs.read(String?),
          )
        )
      end
    rescue ex: DB::NoResultsError
      Err.not_found(ex)
    rescue ex: DB::Error
      Err.fail(ex)
    end

    def find_all(account_id)
      deploys = Array(Deploy).new
      @db.query "SELECT id, account_id, server, local_dir, remote_dir FROM deploys WHERE account_id = ?", account_id do |rs|
        rs.each do
          deploys << Deploy.new(
            id: rs.read(Int64),
            account_id: rs.read(Int64),
            server: rs.read(String?),
            local_dir: rs.read(String?),
            remote_dir: rs.read(String?),
          )
        end
      end

      Ok.done(deploys)
    rescue ex: DB::Error
      Err.fail(ex)
    end

    def create(account_id, params)
      key = Random::Secure.random_bytes(64)
      username = params["username"]?
      password = params["password"]?
      username_iv = username ? Random::Secure.random_bytes(32) : nil
      password_iv = password ? Random::Secure.random_bytes(32) : nil

      result = @db.exec "INSERT INTO deploys (account_id, server, local_dir, remote_dir, username, username_iv, password, password_iv) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
        account_id,
        params["server"]?,
        params["local_dir"]?,
        params["remote_dir"]?,
        encrypt(username, key, username_iv),
        username_iv,
        encrypt(password, key, password_iv),
        password_iv

      Ok.created(result.last_insert_id)
    rescue ex: DB::Error
      Err.fail(ex)
    end

    private def encrypt(data, key, iv)
      return nil if data.nil?

      cipher = OpenSSL::Cipher.new("aes-256-cbc")
      cipher.encrypt
      cipher.key = key
      cipher.random_iv

      io = IO::Memory.new
      io.write(cipher.update(data))
      io.write(cipher.final)
      io.rewind

      Base64.encode(io.to_s)
    end
  end
end


