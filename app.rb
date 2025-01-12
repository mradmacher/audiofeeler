# frozen_string_literal: true

require "sinatra"

get "/" do
  @accounts = []
  erb :"accounts.html", layout: :"application.html"
end

get "/:name" do
  @account = Struct.new(:name, :title, :url).new(params[:name], params[:title], "example.org")
  erb :"account.html", layout: :"application.html"
end
