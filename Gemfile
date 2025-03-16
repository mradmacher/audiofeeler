# frozen_string_literal: true

source "https://rubygems.org"

gem "entitainer", "~> 0.0.1"
gem "optiomist", "~> 0.0.3"
gem "param_param", "~> 1.0.0"
gem "puma"
gem "rackup"
gem "rake"
gem "resonad"
gem "sqlite3"
gem "sinatra", ">= 4.1.1"

gem "jekyll", "~> 4.4.1"
group :jekyll_plugins do
  gem "jekyll-feed"
  gem "jekyll-seo-tag"
end

group :development do
  gem "sinatra-contrib"
end

group :test, :development do
  gem "rubocop"
  gem "rubocop-minitest"
  gem "rubocop-rake"
  gem "rerun"
end

group :test do
  gem "minitest"
  gem "minitest-hooks"
  gem "minitest-rg"
  gem "rack-test"
end

