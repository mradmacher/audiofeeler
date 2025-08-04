require "yaml"

class Event
  def self.new(ctx : YAML::PullParser, node : YAML::Nodes::Node)
    puts node
  end
end

yaml = File.open("data/accounts/czarnymotyl.art/_data/events.yml") do |file|
    YAML.parse(file)
end
yaml.as_a.each do |data|
  puts data
end
