# frozen_string_literal: true

require "entitainer"
require "resonad"

Result = Resonad

module Audiofeeler
  class Account
    include Entitainer

    schema do
      attributes :name

      has_many :events
    end
  end

  class Event
    include Entitainer

    schema do
      attributes :date
      attributes :hour
      attributes :venue
      attributes :place
      attributes :city
      attributes :address
    end
  end
end
