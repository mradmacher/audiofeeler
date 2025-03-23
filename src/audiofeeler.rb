# frozen_string_literal: true

require "entitainer"
require "resonad"

Result = Resonad

module Audiofeeler
  class Account
    include Entitainer

    schema do
      attributes :name

      has_many :events
    end
  end

  class Event
    include Entitainer

    schema do
      attributes :date
      attributes :hour
      attributes :venue
      attributes :place
      attributes :city
      attributes :address
    end
  end

  class AccountRepo
    def initialize(db)
      @db = db
    end

    def create(params)
      @db.execute "INSERT INTO accounts (name) VALUES (?)", [params[:name]]
      Result.success
    end

    def find_all
      result = @db.execute "SELECT id, name FROM accounts"

      Result.success(
        result.map { |r| Audiofeeler::Account.new(id: r[0], name: r[1]) }
      )
    end

    def find_one(id)
      result = @db.execute "SELECT id, name FROM accounts WHERE id = ?", id
      return Result.failure(:not_found) if result.empty?

      Result.success(Audiofeeler::Account.new(id: result.first[0], name: result.first[1]))
    end
  end

  class EventRepo
    def initialize(db)
      @db = db
    end

    def find_all(account_id)
      result = @db.execute "SELECT id, date, hour, venue, place, city, address FROM events WHERE account_id = ?", account_id

      events = result.map do |r|
        Audiofeeler::Event.new(id: r[0], date: r[1], hour: r[2], venue: r[3], place: r[4], city: r[5], address: r[6])
      end

      Result.success(events)
    end

    def find_one(account_id, event_id)
      result = @db.execute "SELECT id, date, hour, venue, place, city, address FROM events WHERE account_id = ? and id = ?", [account_id, event_id]
      return Result.failure(:not_found) if result.empty?

      r = result.first
      Result.success(
        Audiofeeler::Event.new(id: r[0], date: r[1], hour: r[2], venue: r[3], place: r[4], city: r[5], address: r[6])
      )
    end

    def create(account_id, params)
      @db.execute "INSERT INTO events (account_id, date, hour, venue, place, city, address) VALUES (?, ?, ?, ?, ?, ?, ?)",
        [account_id, params[:date], params[:hour], params[:venue], params[:place], params[:city], params[:address]]

      Result.success
    end

    def update(account_id, event_id, params)
      @db.execute "UPDATE events SET date = ?, hour = ?, venue = ?, place = ?, city = ?, address = ? WHERE account_id = ? and id = ?",
        [params[:date], params[:hour], params[:venue], params[:place], params[:city], params[:address], account_id, event_id]

      Result.success
    end
  end
end
