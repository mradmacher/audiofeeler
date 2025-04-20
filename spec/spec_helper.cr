require "spec"
require "sqlite3"
require "../src/*"

TESTDB = DB.open "sqlite3://./data/test.db"

schema = File.read("db/schema.sql")
schema.split(";").each do |stmt|
  stmt = stmt.strip
  next if stmt.empty?

  TESTDB.exec stmt
end
