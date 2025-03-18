# frozen_string_literal: true

$LOAD_PATH << File.join(__dir__, 'lib')

require "sinatra/base"
require "audiofeeler"

require "sqlite3"

DB = SQLite3::Database.new "./data/development.db"

def select_layout(request, layouts)
  if request.env["HTTP_HX_REQUEST"]
    erb layouts.last, layout: false
  else
    erb layouts.first, layout: true do
      erb layouts.last, layout: false
    end
  end
end

def handle_result(result)
  case result
  in { value: }
    yield value
  in { error: :not_found }
    halt 404
  in { error: }
    halt 500
  end
end

# AccountRepo.new(DB).tap do |r|
#   r.create(name: "Czarny motyl")
#   r.create(name: "Iglika")
#   r.create(name: "Etnorozrabiaka")
# end
# EventRepo.new(DB).create(3, { date: "28.03.2025", hour: "19:00", venue: "Festiwal rozrabiaków", place: "Bar rozrabiaków", city: "Rozrabiaków", address: "Rozrabiacka 23" })

module Web
  class Events < Sinatra::Base
    get "/accounts/:account_id/events" do
      result = Audiofeeler::AccountRepo.new(DB).find_one(params[:account_id])
        .on_success {
          @account = it
        }.and_then {
          Audiofeeler::EventRepo.new(DB).find_all(@account.id)
        }

      handle_result(result) do |events|
        @events = events
        select_layout(request, [:account_layout, :events])
      end
    end

    get "/accounts/:account_id/events/new" do
      result = Audiofeeler::AccountRepo.new(DB).find_one(params[:account_id])
        .on_success {
          @account = it
        }

      handle_result(result) do
        @event = Audiofeeler::Event.new
        erb :event_form, layout: false
      end
    end

    get "/accounts/:account_id/events/:id/edit" do
      result = Audiofeeler::AccountRepo.new(DB).find_one(params[:account_id])
        .on_success { @account = it }
        .and_then { Audiofeeler::EventRepo.new(DB).find_one(@account.id, params[:id]) }

      handle_result(result) do |event|
        @event = event
        erb :event_form, layout: false
      end
    end

    post "/accounts/:account_id/events" do
      result = Audiofeeler::AccountRepo.new(DB).find_one(params[:account_id])
        .on_success { @account = it }
        .and_then { Audiofeeler::EventRepo.new(DB).create(@account.id, params[:event]) }

      handle_result(result) do
        redirect "/accounts/#{@account.id}/events", 303
      end
    end

    put "/accounts/:account_id/events/:id" do
      result = Audiofeeler::AccountRepo.new(DB).find_one(params[:account_id])
        .on_success { @account = it }
        .and_then { Audiofeeler::EventRepo.new(DB).update(@account.id, params[:id], params[:event]) }


      handle_result(result) do
        redirect "/accounts/#{@account.id}/events", 303
      end
    end
  end
end

class AudiofeelerWeb < Sinatra::Base
  use Web::Events

  configure :production, :development do
    enable :logging
  end

  get "/" do
    redirect "/accounts", 303
  end

  get "/accounts" do
    result = Audiofeeler::AccountRepo.new(DB).find_all

    handle_result(result) do |accounts|
      @accounts = accounts
      erb :accounts, layout: true
    end
  end

  get "/accounts/:id" do
    result = Audiofeeler::AccountRepo.new(DB).find_one(params[:id])
    handle_result(result) do |account|
      @account = account
      erb :account_layout, layout: true do
        erb :account
      end
    end
  end

  get "/accounts/:account_id/pages" do
    erb :pages, layout: false
  end

  get "/accounts/:account_id/videos" do
    erb :videos, layout: false
  end

  get "/accounts/:account_id/lyrics" do
    erb :lyrics, layout: false
  end
end
