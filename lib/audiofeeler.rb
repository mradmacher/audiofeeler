# frozen_string_literal: true

require 'entitainer'

module Audiofeeler
  class Account
    include Entitainer

    schema do
      attributes :name
    end
  end
end
