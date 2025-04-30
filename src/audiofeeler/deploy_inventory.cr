require "result"
require "openssl"
require "random/secure"

module Audiofeeler
  class DeployInventory
    CIPHER_NAME = "aes-256-cbc"

    def initialize(@db : DB::Database, encryption_key : String)
      @encryption_key = Base64.decode(encryption_key)
    end

    def find_one(account_id, id)
      @db.query_one "SELECT id, account_id, server, local_dir, remote_dir, username, password FROM deploys WHERE account_id = ? and id = ?", account_id, id do |rs|
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
      @db.query "SELECT id, account_id, server, local_dir, remote_dir, username, password FROM deploys WHERE account_id = ?", account_id do |rs|
        rs.each do
          deploys << Deploy.new(
            id: rs.read(Int64),
            account_id: rs.read(Int64),
            server: rs.read(String?),
            local_dir: rs.read(String?),
            remote_dir: rs.read(String?),
            username: rs.read(String?),
            password: rs.read(String?),
          )
        end
      end

      Ok.done(deploys)
    rescue ex: DB::Error
      Err.fail(ex)
    end

    def update(account_id, id, params)
      result = if params.has_key?("server") || params.has_key?("local_dir") || params.has_key?("remote_dir")
        update_paths(account_id, id, params)
      end
      return result.not_nil! if result && result.err?

      result = if params.has_key?("username") || params.has_key?("password")
        update_credentials(account_id, id, params)
      end
      return result.not_nil! if result && result.err?

      Ok.updated(id)
    end

    private def update_paths(account_id, id, params)
      @db.exec "UPDATE deploys SET server = ?, local_dir = ?, remote_dir = ? WHERE account_id = ? AND id = ?", params["server"], params["local_dir"], params["remote_dir"], account_id, id

      Ok.updated(id)
    rescue ex: DB::Error
      Err.fail(ex)
    end

    private def update_credentials(account_id, id, params)
      username = params["username"]
      password = params["password"]
      username_iv = random_iv
      password_iv = random_iv

      @db.exec "UPDATE deploys SET username = ?, username_iv = ?, password = ?, password_iv = ? WHERE account_id = ? AND id = ?",
        encrypt(username, @encryption_key, username_iv),
        username_iv,
        encrypt(password, @encryption_key, password_iv),
        password_iv,
        account_id,
        id

      Ok.updated(id)
    rescue ex: DB::Error
      Err.fail(ex)
    end

    def create(account_id, params)
      username = params["username"]?
      password = params["password"]?
      username_iv = username ? random_iv : nil
      password_iv = password ? random_iv : nil

      result = @db.exec "INSERT INTO deploys (account_id, server, local_dir, remote_dir, username, username_iv, password, password_iv) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
        account_id,
        params["server"]?,
        params["local_dir"]?,
        params["remote_dir"]?,
        encrypt(username, @encryption_key, username_iv),
        username_iv,
        encrypt(password, @encryption_key, password_iv),
        password_iv

      Ok.created(result.last_insert_id)
    rescue ex: DB::Error
      Err.fail(ex)
    end

    def delete(account_id, id)
      result = @db.exec "DELETE FROM deploys WHERE account_id = ? and id = ?", account_id, id

      Ok.destroyed(id)
    rescue ex: DB::Error
      Err.fail(ex)
    end

    def self.random_encryption_key
      Random::Secure.base64(64)
    end

    private def random_iv
      Random::Secure.random_bytes(32)
    end

    private def encrypt(data, key, iv)
      return nil if data.nil?

      cipher = OpenSSL::Cipher.new(CIPHER_NAME)
      cipher.encrypt
      cipher.key = key
      cipher.iv = iv.not_nil!

      io = IO::Memory.new
      io.write(cipher.update(data))
      io.write(cipher.final)
      io.rewind

      Base64.encode(io.to_s)
    end
  end
end


