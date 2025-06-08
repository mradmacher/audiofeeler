require "./spec_helper"

describe "DeploymentInventory" do
  account_inventory = Audiofeeler::AccountInventory.new(TESTDB)
  deployment_inventory = Audiofeeler::DeploymentInventory.new(TESTDB, Audiofeeler::DeploymentInventory.random_encryption_key)

  account = account_inventory.find_one(
    account_inventory.create({"name" => "Test", "source_dir" => "here"}).unwrap
  ).unwrap

  describe "#find_one_decrypted" do
    plain_username = "look at me"
    plain_password = "you can see me"
    result = deployment_inventory.create(account.id, {
      "server" => "example.com",
      "remote_dir" => "there",
      "username" => plain_username,
      "password" => plain_password,
    })
    result.ok?.should be_true

    deploy = deployment_inventory.find_one_decrypted(account.id, result.unwrap).unwrap

    deploy.username.should eq plain_username
    deploy.password.should eq plain_password
  end

  describe "#create" do
    it "creates new deploy" do
      result = deployment_inventory.create(account.id, {
        "server" => "example.com",
        "remote_dir" => "there",
      })
      result.ok?.should be_true

      deploy = deployment_inventory.find_one(account.id, result.unwrap).unwrap

      deploy.server.should eq "example.com"
      deploy.remote_dir.should eq "there"
      deploy.username.should be_nil
      deploy.password.should be_nil
    end

    it "encrypts username and password" do
      plain_username = "look at me"
      plain_password = "you can see me"
      result = deployment_inventory.create(account.id, {
        "server" => "example.com",
        "username" => plain_username,
        "password" => plain_password,
      })
      result.ok?.should be_true

      deploy = deployment_inventory.find_one(account.id, result.unwrap).unwrap

      deploy.username.should_not be_nil
      deploy.username.should_not eq plain_username

      deploy.password.should_not be_nil
      deploy.password.should_not eq plain_password
    end
  end
end
