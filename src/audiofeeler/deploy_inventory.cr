require "result"

module Audiofeeler
  class DeployInventory
    def initialize(@db : DB::Database)
      @db = db
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
      result = @db.exec "INSERT INTO deploys (account_id, server, local_dir, remote_dir, username, password) VALUES (?, ?, ?, ?, ?, ?)",
        account_id, params["deployl[server]"], params["deploy[local_dir]"], params["deploy[remote_dir]"], params["deploy[username]"], params["deploy[password]"]

      Ok.created(result.last_insert_id)
    rescue ex: DB::Error
      Err.fail(ex)
    end
  end
end


