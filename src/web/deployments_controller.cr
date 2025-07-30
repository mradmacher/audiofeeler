module Web
  module DeploymentsController
    def self.configure_routes(app)
      accounts_inventory = Audiofeeler::AccountInventory.new(app.db)
      deployment_inventory = Audiofeeler::DeploymentInventory.new(app.db, app.encryption_key)

      app.router.add_route "GET", "/accounts/:id/deployments/new", do |env|
        result = accounts_inventory.find_one(env.params.url["id"])
        handle_result(result, env) do |account|
          deployment = Audiofeeler::Deployment.new
          render_no_layout("deployment_form")
        end
      end

      app.router.add_route "GET", "/accounts/:id/deployments/:deployment_id/edit", do |env|
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

      app.router.add_route "POST", "/accounts/:id/deployments", do |env|
        result = accounts_inventory.find_one(env.params.url["id"])
        handle_result(result, env) do |account|
          result = deployment_inventory.create(account.id, env.params.body)
          handle_result(result, env) do
            env.redirect "/accounts/#{account.id}/config", 303
          end
        end
      end

      app.router.add_route "PUT", "/accounts/:id/deployments/:deployment_id", do |env|
        result = accounts_inventory.find_one(env.params.url["id"])
        handle_result(result, env) do |account|
          result = deployment_inventory.update(account.id, env.params.url["deployment_id"], env.params.body)
          handle_result(result, env) do
            env.redirect "/accounts/#{account.id}/config", 303
          end
        end
      end

      app.router.add_route "POST", "/accounts/:id/deployments/:deployment_id/release", do |env|
        result = accounts_inventory.find_one(env.params.url["id"])
        handle_result(result, env) do |account|
          result = deployment_inventory.find_one_decrypted(account.id, env.params.url["deployment_id"])
          handle_result(result, env) do |deployment|
            stdout = IO::Memory.new
            build_status = Process.run("bundle", ["exec", "jekyll", "build", "--incremental", "-s", "data/accounts/#{account.source_dir}/", "-d", "data/accounts/#{account.source_dir}/_site/"], output: stdout)
            build_output = stdout.to_s

            stdout = IO::Memory.new
            deployment_status = Process.run("npx", ["ftp-deploy", "--server", deployment.server.not_nil!, "--username", deployment.username.not_nil!, "--password", deployment.password.not_nil!, "--local-dir", "data/accounts/#{account.source_dir}/_site/", "--server-dir", deployment.remote_dir.not_nil!], output: stdout)
            deployment_output = stdout.to_s
            render_no_layout("deploy_result")
          end
        end
      end

      app.router.add_route "DELETE", "/accounts/:id/deployments/:deployment_id", do |env|
        result = accounts_inventory.find_one(env.params.url["id"])
        handle_result(result, env) do |account|
          result = deployment_inventory.delete(account.id, env.params.url["deployment_id"])
          handle_result(result, env) do
            env.redirect "/accounts/#{account.id}/config", 303
          end
        end
      end
    end
  end
end
