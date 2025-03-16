# frozen_string_literal: true

$LOAD_PATH << File.join(__dir__, 'lib')

require "sinatra/base"
require "audiofeeler"

require "sqlite3"

DB = SQLite3::Database.new ":memory:"


DB.execute <<SQL
  CREATE TABLE accounts (
    id INTEGER PRIMARY KEY,
    name TEXT,
    dir TEXT
  );
SQL

DB.execute <<SQL
  CREATE TABLE deploys (
    id INTEGER PRIMARY KEY,
    account_id INTEGER,
    server TEXT,
    username TEXT,
    username_iv TEXT,
    password TEXT,
    password_iv TEXT,
    dir TEXT,
    FOREIGN KEY(account_id) REFERENCES accounts(id)
  );
SQL

DB.execute <<SQL
  CREATE TABLE events (
    id INTEGER PRIMARY KEY,
    account_id INTEGER,
    date TEXT,
    hour TEXT,
    venue TEXT,
    place TEXT,
    address TEXT,
    FOREIGN KEY(account_id) REFERENCES accounts(id)
  );
SQL

class AccountRepo
  def initialize(db)
    @db = db
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
    result = @db.execute "SELECT id, date, hour, venue, place, address FROM events WHERE account_id = ?", account_id

    events = result.map do |r|
      Audiofeeler::Event.new(id: r[0], date: r[1], hour: r[2], venue: r[3], place: r[4], address: r[5])
    end

    Result.success(events)
  end

  def find_one(account_id, event_id)
    result = @db.execute "SELECT id, date, hour, venue, place, address FROM events WHERE account_id = ? and id = ?", [account_id, event_id]
    return Result.failure(:not_found) if result.empty?

    r = result.first
    Result.success(
      Audiofeeler::Event.new(id: r[0], date: r[1], hour: r[2], venue: r[3], place: r[4], address: r[5])
    )
  end

  def create(account_id, params)
    @db.execute "INSERT INTO events (account_id, date, hour, venue, place, address) VALUES (?, ?, ?, ?, ?, ?)",
      [account_id, params[:date], params[:hour], params[:venue], params[:place], params[:address]]

    Result.success
  end

  def update(account_id, event_id, params)
    @db.execute "UPDATE events SET date = ?, hour = ?, venue = ?, place = ?, address = ? WHERE account_id = ? and id = ?",
      [params[:date], params[:hour], params[:venue], params[:place], params[:address], account_id, event_id]

    Result.success
  end
end

def select_layout(request, layouts)
  if request.env["HTTP_HX_REQUEST"]
    erb layouts.last, layout: false
  else
    erb layouts.first, layout: true do
      erb layouts.last, layout: false
    end
  end
end

def handle_result(result)
  case result
  in { value: }
    yield value
  in { error: :not_found }
    halt 404
  in { error: }
    halt 500
  end
end

DB.execute "INSERT INTO accounts (name) VALUES (?)", "Czarny motyl"
DB.execute "INSERT INTO accounts (name) VALUES (?)", "Iglika"
DB.execute "INSERT INTO accounts (name) VALUES (?)", "Etnorozrabiaka"
EventRepo.new(DB).create(3, { date: "28.03.2025", hour: "19:00", venue: "Festiwal rozrabiaków", place: "Bar rozrabiaków", address: "Rozrabiacka 23, Rozrabiaków" })

module Web
  class Events < Sinatra::Base
    get "/accounts/:account_id/events" do
      result = AccountRepo.new(DB).find_one(params[:account_id])
        .on_success {
          @account = it
        }.and_then {
          EventRepo.new(DB).find_all(@account.id)
        }

      handle_result(result) do |events|
        @events = events
        select_layout(request, [:account_layout, :events])
      end
    end

    get "/accounts/:account_id/events/new" do
      result = AccountRepo.new(DB).find_one(params[:account_id])
        .on_success {
          @account = it
        }

      handle_result(result) do
        @event = Audiofeeler::Event.new
        erb :event_form, layout: false
      end
    end

    get "/accounts/:account_id/events/:id/edit" do
      result = AccountRepo.new(DB).find_one(params[:account_id])
        .on_success { @account = it }
        .and_then { EventRepo.new(DB).find_one(@account.id, params[:id]) }

      handle_result(result) do |event|
        @event = event
        erb :event_form, layout: false
      end
    end

    post "/accounts/:account_id/events" do
      result = AccountRepo.new(DB).find_one(params[:account_id])
        .on_success { @account = it }
        .and_then { EventRepo.new(DB).create(@account.id, params[:event]) }

      handle_result(result) do
        redirect "/accounts/#{@account.id}/events", 303
      end
    end

    put "/accounts/:account_id/events/:id" do
      result = AccountRepo.new(DB).find_one(params[:account_id])
        .on_success { @account = it }
        .and_then { EventRepo.new(DB).update(@account.id, params[:id], params[:event]) }


      handle_result(result) do
        redirect "/accounts/#{@account.id}/events", 303
      end
    end
  end
end

class AudiofeelerWeb < Sinatra::Base
  use Web::Events

  configure :production, :development do
    enable :logging
  end

  get "/" do
    redirect "/accounts", 303
  end

  get "/accounts" do
    result = AccountRepo.new(DB).find_all

    handle_result(result) do |accounts|
      @accounts = accounts
      erb :accounts, layout: true
    end
  end

  get "/accounts/:id" do
    result = AccountRepo.new(DB).find_one(params[:id])
    handle_result(result) do |account|
      @account = account
      erb :account_layout, layout: true do
        erb :account
      end
    end
  end

  get "/accounts/:account_id/pages" do
    erb :pages, layout: false
  end

  get "/accounts/:account_id/videos" do
    erb :videos, layout: false
  end

  get "/accounts/:account_id/lyrics" do
    erb :lyrics, layout: false
  end
end
