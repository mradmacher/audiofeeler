require "result"

module Audiofeeler
  class Account
    getter id, name

    def initialize(
      @id : Int64?,
      @name : String?
    )
    end
  end

  class Event
    getter id, date, hour, venue, place, city, address

    def initialize(
      @id : Int64? = nil,
      @date : String? = nil,
      @hour : String? = nil,
      @venue : String? = nil,
      @place : String? = nil,
      @city : String? = nil,
      @address : String? = nil
    )
    end
  end

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

  class EventInventory
    def initialize(@db : DB::Database)
      @db = db
    end

    def find_all(account_id)
      events = Array(Event).new
      @db.query "SELECT id, date, hour, venue, place, city, address FROM events WHERE account_id = ?", account_id do |rs|
        rs.each do
          events << Event.new(
            id: rs.read(Int64),
            date: rs.read(String?),
            hour: rs.read(String?),
            venue: rs.read(String?),
            place: rs.read(String?),
            city: rs.read(String?),
            address: rs.read(String?),
          )
        end
      end

      Ok.done(events)
    rescue ex: DB::Error
      Err.fail(ex)
    end

    def find_one(account_id, event_id)
      @db.query_one "SELECT id, date, hour, venue, place, city, address FROM events WHERE account_id = ? and id = ?", account_id, event_id do |rs|
        return Ok.done(
          Event.new(
            id: rs.read(Int64),
            date: rs.read(String?),
            hour: rs.read(String?),
            venue: rs.read(String?),
            place: rs.read(String?),
            city: rs.read(String?),
            address: rs.read(String?),
          )
        )
      end
    rescue ex: DB::NoResultsError
      Err.not_found(ex)
    rescue ex: DB::Error
      Err.fail(ex)
    end

    def create(account_id, params)
      result = @db.exec "INSERT INTO events (account_id, date, hour, venue, place, city, address) VALUES (?, ?, ?, ?, ?, ?, ?)",
        account_id, params["event[date]"], params["event[hour]"], params["event[venue]"], params["event[place]"], params["event[city]"], params["event[address]"]

      Ok.created(result.last_insert_id)
    rescue ex: DB::Error
      Err.fail(ex)
    end

    def update(account_id, event_id, params)
      result = @db.exec "UPDATE events SET date = ?, hour = ?, venue = ?, place = ?, city = ?, address = ? WHERE account_id = ? and id = ?",
        params["event[date]"]?, params["event[hour]"]?, params["event[venue]"]?, params["event[place]"]?, params["event[city]"]?, params["event[address]"], account_id, event_id

      Ok.updated(event_id)
    rescue ex: DB::Error
      Err.fail(ex)
    end
  end
end
