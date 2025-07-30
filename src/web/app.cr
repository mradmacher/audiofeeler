require "kemal"
require "sqlite3"
require "../audiofeeler"
require "./accounts_controller"
require "./events_controller"
require "./deployments_controller"


def handle_result(result, env)
  case result
  when Ok
    yield result.unwrap
  when .status? :not_found
    env.response.status_code = 404
    env.response.close
  when Err
    Log.error(result.unwrap)
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
  render "src/web/views/#{ {{filename}} }.ecr"
end

macro render_with_layout(filename)
  render "src/web/views/#{ {{filename}} }.ecr", "src/web/views/layout.ecr"
end

macro render_htmx(xhr, filename)
  {{xhr}} ? render_no_layout({{filename}}) : render_with_layout({{filename}})
end

class App
  getter db : DB::Database
  getter :encryption_key
  getter router : Kemal::RouteHandler

  def initialize
    @encryption_key = ENV["AUDIOFEELER_ENCRYPTION_KEY"]
    @db = DB.open "sqlite3://./data/development.db"
    @router = Kemal::RouteHandler::INSTANCE
  end
end

app = App.new

accounts_inventory = Audiofeeler::AccountInventory.new(app.db)
deployment_inventory = Audiofeeler::DeploymentInventory.new(app.db, app.encryption_key)

Web::AccountsController.configure_routes(app)
Web::EventsController.configure_routes(app)
Web::DeploymentsController.configure_routes(app)

get "/" do |env|
  env.redirect "/accounts", 303
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
