require "./spec_helper"

describe "DeployInventory" do
  account_inventory = Audiofeeler::AccountInventory.new(TESTDB)

  account = account_inventory.find_one(
    account_inventory.create({"name" => "Test"}).unwrap
  ).unwrap

  it "works" do
    account.name.should eq("Test")
  end
end
