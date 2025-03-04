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
    venue TEXT,
    FOREIGN KEY(account_id) REFERENCES accounts(id)
  );
SQL

def find_account(id)
  result = DB.execute "SELECT id, name FROM accounts WHERE id = ?", id
  return nil if result.empty?

  Audiofeeler::Account.new(id: result.first[0], name: result.first[1])
end

def find_account_events(account_id)
  result = DB.execute "SELECT id, venue FROM events WHERE account_id = ?", account_id

  result.map do |r|
    Audiofeeler::Event.new(id: r[0], venue: r[1])
  end
end

def add_account_event(account_id, params)
  DB.execute "INSERT INTO events (account_id, venue) VALUES (?, ?)", [account_id, params[:venue]]
end

def no_layout_or(request, layout)
  request.env["HTTP_HX_REQUEST"] ? false : layout
end

DB.execute "INSERT INTO accounts (name) VALUES (?)", "Czarny motyl"
DB.execute "INSERT INTO accounts (name) VALUES (?)", "Iglika"
DB.execute "INSERT INTO accounts (name) VALUES (?)", "Etnorozrabiaka"
add_account_event(3, { venue: "Festiwal rozrabiaków" })

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

get "/accounts/:id/events" do
  @account = find_account(params[:id])
  halt 404 if @account.nil?

  @events = find_account_events(@account.id)
  if request.env["HTTP_HX_REQUEST"]
    erb :"events.html", layout: false
  else
    erb(:account_layout, layout: true) do
      erb :"events.html"
    end
  end
end

post "/accounts/:id/events" do
  @account = find_account(params[:id])
  halt 404 if @account.nil?

  add_account_event(@account.id, params[:event])

  redirect "/accounts/#{@account.id}/events", 303
end

get "/accounts/:id/pages" do
  erb :"pages.html", layout: false
end

get "/accounts/:id/videos" do
  erb :"videos.html", layout: false
end

get "/accounts/:id/lyrics" do
  erb :"lyrics.html", layout: false
end
