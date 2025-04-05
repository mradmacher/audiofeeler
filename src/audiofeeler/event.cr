module Audiofeeler
  struct Event
    getter id, date, hour, venue, place, city, address

    def initialize(
      @id : Int64? = nil,
      @date : String? = nil,
      @hour : String? = nil,
      @venue : String? = nil,
      @place : String? = nil,
      @city : String? = nil,
      @address : String? = nil
    )
    end
  end
end
