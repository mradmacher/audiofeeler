require "result"

module Audiofeeler
  class AccountInventory
    def initialize(@db : DB::Database)
      @db = db
    end

    def create(params)
      er = @db.exec "INSERT INTO accounts (name) VALUES (?)", params[:name]
      Ok.created(er.last_insert_id)
    rescue ex: DB::Error
      Err.fail(ex)
    end

    def find_all
      accounts = Array(Account).new
      rs = @db.query "SELECT id, name FROM accounts"

      rs.each do
        accounts << Account.new(id: rs.read(Int64), name: rs.read(String?))
      end
      Ok.done(accounts)
    rescue ex: DB::Error
      Err.fail(ex)
    end

    def find_one(id)
      @db.query_one "SELECT id, name FROM accounts WHERE id = ?", id do |rs|
        return Ok.done(Account.new(id: rs.read(Int64), name: rs.read(String)))
      end
    rescue ex: DB::NoResultsError
      Err.not_found(ex)
    rescue ex: DB::Error
      Err.fail(ex)
    end
  end
end
