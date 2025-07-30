require "./spec_helper"
require "spec-kemal"
require "../web"

describe "Web" do
  it "renders /" do
    get "/"

    response.body.should eq 200
  end
end
