module Web
  module EventsController
    def self.configure_routes(app)
      accounts_inventory = Audiofeeler::AccountInventory.new(app.db)
      events_inventory = Audiofeeler::EventInventory.new(app.db)

      app.router.add_route "GET", "/accounts/:id/events", do |env|
        result = accounts_inventory.find_one(env.params.url["id"])
        handle_result(result, env) do |account|
          result = events_inventory.find_all(account.id)
          handle_result(result, env) do |events|
            render_htmx(is_xhr(env), "events")
          end
        end
      end

      app.router.add_route "GET", "/accounts/:id/events/new", do |env|
        result = accounts_inventory.find_one(env.params.url["id"])
        handle_result(result, env) do |account|
          event = Audiofeeler::Event.new
          render_no_layout("event_form")
        end
      end

      app.router.add_route "GET", "/accounts/:id/events/:eid/edit", do |env|
        result = accounts_inventory.find_one(env.params.url["id"])
        handle_result(result, env) do |account|
          result = events_inventory.find_one(account.id, env.params.url["eid"])
          handle_result(result, env) do |event|
            render_no_layout("event_form")
          end
        end
      end

      app.router.add_route "POST", "/accounts/:id/events", do |env|
        result = accounts_inventory.find_one(env.params.url["id"])
        handle_result(result, env) do |account|
          result = events_inventory.create(account.id, env.params.body)
          handle_result(result, env) do
            env.redirect "/accounts/#{account.id}/events", 303
          end
        end
      end

      app.router.add_route "PUT", "/accounts/:id/events/:eid", do |env|
        result = accounts_inventory.find_one(env.params.url["id"])
        handle_result(result, env) do |account|
          result = events_inventory.update(account.id, env.params.url["eid"], env.params.body)
          handle_result(result, env) do
            env.redirect "/accounts/#{account.id}/events", 303
          end
        end
      end

      app.router.add_route "DELETE", "/accounts/:id/events/:eid", do |env|
        result = accounts_inventory.find_one(env.params.url["id"])
        handle_result(result, env) do |account|
          result = events_inventory.delete(account.id, env.params.url["eid"])
          handle_result(result, env) do
            env.redirect "/accounts/#{account.id}/events", 303
          end
        end
      end
    end
  end
end
