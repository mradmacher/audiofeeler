# frozen_string_literal: true

$LOAD_PATH << File.join(__dir__, 'lib')

require "sinatra"
require "audiofeeler"

require "sqlite3"

DB = SQLite3::Database.new ":memory:"

DB.execute <<SQL
  CREATE TABLE accounts (
    id INTEGER PRIMARY KEY,
    name TEXT
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

def find_account(id)
  result = DB.execute "SELECT id, name FROM accounts WHERE id = ?", id
  return nil if result.empty?

  Audiofeeler::Account.new(id: result.first[0], name: result.first[1])
end

def find_account_events(account_id)
  result = DB.execute "SELECT id, date, hour, venue, place, address FROM events WHERE account_id = ?", account_id

  result.map do |r|
    Audiofeeler::Event.new(id: r[0], date: r[1], hour: r[2], venue: r[3], place: r[4], address: r[5])
  end
end

def find_account_event(account_id, event_id)
  result = DB.execute "SELECT id, date, hour, venue, place, address FROM events WHERE account_id = ? and id = ?", [account_id, event_id]
  return nil if result.empty?

  r = result.first

  Audiofeeler::Event.new(id: r[0], date: r[1], hour: r[2], venue: r[3], place: r[4], address: r[5])
end

def add_account_event(account_id, params)
  DB.execute "INSERT INTO events (account_id, date, hour, venue, place, address) VALUES (?, ?, ?, ?, ?, ?)",
    [account_id, params[:date], params[:hour], params[:venue], params[:place], params[:address]]
end

def update_account_event(account_id, event_id, params)
  DB.execute "UPDATE events SET date = ?, hour = ?, venue = ?, place = ?, address = ? WHERE account_id = ? and id = ?",
    [params[:date], params[:hour], params[:venue], params[:place], params[:address], account_id, event_id]
end

def no_layout_or(request, layout)
  request.env["HTTP_HX_REQUEST"] ? false : layout
end

DB.execute "INSERT INTO accounts (name) VALUES (?)", "Czarny motyl"
DB.execute "INSERT INTO accounts (name) VALUES (?)", "Iglika"
DB.execute "INSERT INTO accounts (name) VALUES (?)", "Etnorozrabiaka"
add_account_event(3, { date: "28.03.2025", hour: "19:00", venue: "Festiwal rozrabiaków", place: "Bar rozrabiaków", address: "Rozrabiacka 23, Rozrabiaków" })

get "/" do
  redirect "/accounts", 303
end

get "/accounts" do
  result = DB.execute "SELECT id, name FROM accounts"
  @accounts = result.map do |row|
    Audiofeeler::Account.new(id: row[0], name: row[1])
  end

  erb :"accounts.html", layout: true
end

get "/accounts/:id" do
  @account = find_account(params[:id])
  halt 404 if @account.nil?

  erb :account_layout, layout: true do
    erb :account
  end
end

get "/accounts/:account_id/events" do
  @account = find_account(params[:account_id])
  halt 404 if @account.nil?

  @events = find_account_events(@account.id)
  if request.env["HTTP_HX_REQUEST"]
    erb :events, layout: false
  else
    erb(:account_layout, layout: true) do
      erb(:events)
    end
  end
end

get "/accounts/:account_id/events/new" do
  @account = find_account(params[:account_id])
  halt 404 if @account.nil?

  @event = Audiofeeler::Event.new

  erb :event_form, layout: false
end

get "/accounts/:account_id/events/:id/edit" do
  @account = find_account(params[:account_id])
  halt 404 if @account.nil?

  @event = find_account_event(@account.id, params[:id])

  erb :event_form, layout: false
end

post "/accounts/:account_id/events" do
  @account = find_account(params[:account_id])
  halt 404 if @account.nil?

  add_account_event(@account.id, params[:event])

  redirect "/accounts/#{@account.id}/events", 303
end

put "/accounts/:account_id/events/:id" do
  @account = find_account(params[:account_id])
  halt 404 if @account.nil?

  @event = update_account_event(@account.id, params[:id], params[:event])

  redirect "/accounts/#{@account.id}/events", 303
end

get "/accounts/:account_id/pages" do
  erb :"pages.html", layout: false
end

get "/accounts/:account_id/videos" do
  erb :"videos.html", layout: false
end

get "/accounts/:account_id/lyrics" do
  erb :"lyrics.html", layout: false
end
