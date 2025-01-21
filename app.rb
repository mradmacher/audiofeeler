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

DB.execute "INSERT INTO accounts (name) VALUES (?)", "Czarny motyl"
DB.execute "INSERT INTO accounts (name) VALUES (?)", "Iglika"
DB.execute "INSERT INTO accounts (name) VALUES (?)", "Etnorozrabiaka"

get "/" do
  redirect "/accounts", 303
end

get "/accounts" do
  result = DB.execute "SELECT id, name FROM accounts"
  @accounts = result.map do |row|
    Audiofeeler::Account.new(id: row[0], name: row[1])
  end

  erb :"accounts.html", layout: :"application.html"
end

get "/accounts/:id" do
  result = DB.execute "SELECT id, name FROM accounts WHERE id = ?", params[:id]
  halt 404 if result.empty?

  @account = Audiofeeler::Account.new(id: result.first[0], name: result.first[1])
  erb :"account.html", layout: :"application.html"
end

get "/accounts/:id/events" do
  erb :"events.html"
end

get "/accounts/:id/pages" do
  erb :"pages.html"
end

get "/accounts/:id/videos" do
  erb :"videos.html"
end

get "/accounts/:id/lyrics" do
  erb :"lyrics.html"
end
