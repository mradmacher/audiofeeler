require "sqlite3"
require "./src/audiofeeler"

db = DB.open "sqlite3://./data/development.db"

aRepo = Audiofeeler::AccountInventory.new(db)
eRepo = Audiofeeler::EventInventory.new(db)

result = aRepo.find_all
case result
in Ok
  result.unwrap.each do |account|
    puts "#{account.name} (#{account.id})"
  end
in Err
  puts result.value
end

[3, 30].each do |id|
  result = aRepo.find_one(id)
  case result
  in Ok
    puts eRepo.find_all(result.unwrap.id).unwrap
  in Err
    puts result.value
  end
end
