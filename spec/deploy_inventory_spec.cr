require "./spec_helper"

describe "DeployInventory" do
  account_inventory = Audiofeeler::AccountInventory.new(TESTDB)
  deploy_inventory = Audiofeeler::DeployInventory.new(TESTDB)

  account = account_inventory.find_one(
    account_inventory.create({"name" => "Test"}).unwrap
  ).unwrap

  describe "#create" do
    it "creates new deploy" do
      result = deploy_inventory.create(account.id, {
        "server" => "example.com",
        "local_dir" => "here",
        "remote_dir" => "there",
      })
      result.ok?.should be_true

      result = deploy_inventory.find_one(result.unwrap)
      result.ok?.should be_true
      deploy = result.unwrap

      deploy.server.should eq "example.com"
      deploy.local_dir.should eq "here"
      deploy.remote_dir.should eq "there"
      deploy.username.should be_nil
      deploy.password.should be_nil
    end

    it "encrypts username and password" do
      plain_username = "look at me"
      plain_password = "you can see me"
      result = deploy_inventory.create(account.id, {
        "server" => "example.com",
        "username" => plain_username,
        "password" => plain_password,
      })
      result.ok?.should be_true

      result = deploy_inventory.find_one(result.unwrap)
      result.ok?.should be_true
      deploy = result.unwrap

      deploy.username.should_not be_nil
      deploy.username.should_not eq plain_username

      deploy.password.should_not be_nil
      deploy.password.should_not eq plain_password
    end
  end
end
