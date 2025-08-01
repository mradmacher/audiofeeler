require "./spec_helper"
require "spec-kemal"
require "../src/web/app"

describe "Web" do
  it "renders /" do
    get "/"

    response.status.should eq 200
  end
end
