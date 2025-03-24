require "sqlite3"
require "./src/audiofeeler"

db = DB.open "sqlite3://./data/development.db"

aRepo = Audiofeeler::AccountRepo.new(db)
eRepo = Audiofeeler::EventRepo.new(db)

aRepo.find_all.each do |account|
  puts "#{account.name} (#{account.id})"
end
account = aRepo.find_one(3)
if account
  puts eRepo.find_all(account.id)
end
