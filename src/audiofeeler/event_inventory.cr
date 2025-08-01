require "result"

module Audiofeeler
  class EventInventory
    def initialize(@db : DB::Database)
      @db = db
    end

    def find_all(account_id)
      events = Array(Event).new
      @db.query "SELECT id, name, date, hour, venue, town, location, status FROM events WHERE account_id = ?", account_id do |rs|
        rs.each do
          events << Event.new(
            id: rs.read(Int64),
            name: rs.read(String?),
            date: rs.read(String?),
            hour: rs.read(String?),
            venue: rs.read(String?),
            town: rs.read(String?),
            location: rs.read(String?),
            status: EventStatus.new(rs.read(Int32)),
          )
        end
      end

      Ok.done(events)
    rescue ex: DB::Error
      Err.fail(ex)
    end

    def find_one(account_id, event_id)
      @db.query_one "SELECT id, name, date, hour, venue, town, location, status FROM events WHERE account_id = ? and id = ?", account_id, event_id do |rs|
        return Ok.done(
          Event.new(
            id: rs.read(Int64),
            name: rs.read(String?),
            date: rs.read(String?),
            hour: rs.read(String?),
            venue: rs.read(String?),
            town: rs.read(String?),
            location: rs.read(String?),
            status: EventStatus.new(rs.read(Int32)),
          )
        )
      end
    rescue ex: DB::NoResultsError
      Err.not_found(ex)
    rescue ex: DB::Error
      Err.fail(ex)
    end

    def create(account_id, params)
      result = @db.exec "INSERT INTO events (account_id, name, date, hour, venue, town, location, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
        account_id, params["event[name]"], params["event[date]"], params["event[hour]"], params["event[venue]"], params["event[town]"], params["event[location]"], params["event[status]"]

      Ok.created(result.last_insert_id)
    rescue ex: DB::Error
      Err.fail(ex)
    end

    def update(account_id, event_id, params)
      result = @db.exec "UPDATE events SET name = ?, date = ?, hour = ?, venue = ?, town = ?, location = ?, status = ? WHERE account_id = ? and id = ?",
        params["event[name]"], params["event[date]"]?, params["event[hour]"]?, params["event[venue]"]?, params["event[town]"]?, params["event[location]"]?, params["event[status]"], account_id, event_id

      Ok.updated(event_id)
    rescue ex: DB::Error
      Err.fail(ex)
    end

    def delete(account_id, event_id)
      result = @db.exec "DELETE FROM events WHERE account_id = ? and id = ?", account_id, event_id

      Ok.destroyed(event_id)
    rescue ex: DB::Error
      Err.fail(ex)
    end
  end
end
