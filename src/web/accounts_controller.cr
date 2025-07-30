module Web
  module AccountsController
    def self.configure_routes(app)
      accounts_inventory = Audiofeeler::AccountInventory.new(app.db)
      deployment_inventory = Audiofeeler::DeploymentInventory.new(app.db, app.encryption_key)

      app.router.add_route "GET", "/accounts", do |env|
        result = accounts_inventory.find_all
        handle_result(result, env) do |accounts|
          account = nil
          render_with_layout "accounts"
        end
      end

      app.router.add_route "GET", "/accounts/:id", do |env|
        result = accounts_inventory.find_one(env.params.url["id"])
        handle_result(result, env) do |account|
          deployments = Array(Audiofeeler::Deployment).new #deployment_inventory.find_all(account.id).unwrap
          render_htmx(is_xhr(env), "account")
        end
      end
    end
  end
end
