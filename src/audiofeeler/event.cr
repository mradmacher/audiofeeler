module Audiofeeler
  enum EventStatus
    Current
    Archived
  end

  struct Event
    getter id, name, date, hour, venue, town, location, status

    def initialize(
      @id : Int64? = nil,
      @name : String? = nil,
      @date : String? = nil,
      @hour : String? = nil,
      @venue : String? = nil,
      @town : String? = nil,
      @location : String? = nil,
      @status : EventStatus = EventStatus::Current
    )
    end
  end
end
