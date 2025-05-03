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

macro render_htmx(xhr, filename)
  {{xhr}} ? render_no_layout({{filename}}) : render_with_layout({{filename}})
end

db = DB.open "sqlite3://./data/development.db"

accounts_inventory = Audiofeeler::AccountInventory.new(db)
events_inventory = Audiofeeler::EventInventory.new(db)
deployment_inventory = Audiofeeler::DeploymentInventory.new(db, Audiofeeler::DeploymentInventory.random_encryption_key)

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
    deployments = deployment_inventory.find_all(account.id).unwrap
    render_htmx(is_xhr(env), "account")
  end
end

get "/accounts/:id/events" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    result = events_inventory.find_all(account.id)
    handle_result(result, env) do |events|
      render_htmx(is_xhr(env), "events")
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

get "/accounts/:id/deployments/new" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    deployment = Audiofeeler::Deployment.new
    render_no_layout("deployment_form")
  end
end

get "/accounts/:id/deployments/:deployment_id/edit" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    result = deployment_inventory.find_one(account.id, env.params.url["deployment_id"])
    handle_result(result, env) do |deployment|
      if env.params.query["view"] == "credentials"
        render_no_layout("deployment_credentials_form")
      else
        render_no_layout("deployment_paths_form")
      end
    end
  end
end

post "/accounts/:id/deployments" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    result = deployment_inventory.create(account.id, env.params.body)
    handle_result(result, env) do
      env.redirect "/accounts/#{account.id}/config", 303
    end
  end
end

put "/accounts/:id/deployments/:deployment_id" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    result = deployment_inventory.update(account.id, env.params.url["deployment_id"], env.params.body)
    handle_result(result, env) do
      env.redirect "/accounts/#{account.id}/config", 303
    end
  end
end

post "/accounts/:id/deployments/:deployment_id" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    result = deployment_inventory.find_one(account.id, env.params.url["deployment_id"])
    handle_result(result, env) do |deploy|
      stdout = IO::Memory.new
      process = Process.new("bundle", ["exec", "jekyll", "build", "--incremental", "-s", "data/accounts/#{account.source_dir}/", "-d", "data/accounts/#{account.source_dir}/_site/"], output: stdout)
      status = process.wait
      output = stdout.to_s
      render_no_layout("deploy_result")
    end
  end
end

delete "/accounts/:id/deployments/:deployment_id" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    result = deployment_inventory.delete(account.id, env.params.url["deployment_id"])
    handle_result(result, env) do
      env.redirect "/accounts/#{account.id}/config", 303
    end
  end
end


get "/accounts/:id/config" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    result2 = deployment_inventory.find_all(account.id)
    handle_result(result2, env) do |deployments|
      render_htmx(is_xhr(env), "config")
    end
  end
end

get "/accounts/:id/pages" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    render_htmx(is_xhr(env), "pages")
  end
end

get "/accounts/:id/videos" do |env|
  result = accounts_inventory.find_one(env.params.url["id"])
  handle_result(result, env) do |account|
    render_htmx(is_xhr(env), "videos")
  end
end

Kemal.run
