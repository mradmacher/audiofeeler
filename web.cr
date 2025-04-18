require "kemal"
require "sqlite3"
require "./src/audiofeeler"

def handle_result(result, env)
  case result
  when Ok
    yield result.unwrap
  when .status? :not_found
    env.response.status_code = 404
    env.response.close
  when Err
    env.response.status_code = 500
    env.response.close
  end
end

def is_xhr(env)
  env.request.headers["HX_REQUEST"]?
end

def handle_render(filename, xhr)
  if env.request.headers["HTTP_HX_REQUEST"]
    display0 filename
  else
    display filename
  end
end

macro render_no_layout(filename)
  render "views/#{ {{filename}} }.ecr"
end

macro render_with_layout(filename)
  render "views/#{ {{filename}} }.ecr", "views/layout.ecr"
end

db = DB.open "sqlite3://./data/development.db"

accounts_inventory = Audiofeeler::AccountInventory.new(db)
events_inventory = Audiofeeler::EventInventory.new(db)

get "/" do |env|
  env.redirect "/accounts", 303
end

get "/accounts" do |env|
  result = accounts_inventory.find_all
  handle_result(result, env) do |accounts|
    account = nil
    render_with_layout "accounts"
  end
end

get "/accounts/:id" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    render_with_layout "account"
  end
end

get "/accounts/:id/events" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    result = events_inventory.find_all(account.id)
    handle_result(result, env) do |events|
      is_xhr(env) ? render_no_layout("events") : render_with_layout("events")
    end
  end
end

get "/accounts/:id/events/new" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    event = Audiofeeler::Event.new
    render_no_layout("event_form")
  end
end

get "/accounts/:id/events/:eid/edit" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    result = events_inventory.find_one(account.id, env.params.url["eid"])
    handle_result(result, env) do |event|
      render_no_layout("event_form")
    end
  end
end

post "/accounts/:id/events" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    result = events_inventory.create(account.id, env.params.body)
    handle_result(result, env) do
      env.redirect "/accounts/#{account.id}/events", 303
    end
  end
end

put "/accounts/:id/events/:eid" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    result = events_inventory.update(account.id, env.params.url["eid"], env.params.body)
    handle_result(result, env) do
      env.redirect "/accounts/#{account.id}/events", 303
    end
  end
end

delete "/accounts/:id/events/:eid" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    result = events_inventory.delete(account.id, env.params.url["eid"])
    handle_result(result, env) do
      env.redirect "/accounts/#{account.id}/events", 303
    end
  end
end

get "/accounts/:id/pages" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    is_xhr(env) ? render_no_layout("pages") : render_with_layout("pages")
  end
end

get "/accounts/:id/videos" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    is_xhr(env) ? render_no_layout("videos") : render_with_layout("videos")
  end
end

get "/accounts/:id/lyrics" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    is_xhr(env) ? render_no_layout("lyrics") : render_with_layout("lyrics")
  end
end

Kemal.run
