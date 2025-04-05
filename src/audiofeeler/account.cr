module Audiofeeler
  struct Account
    getter id, name

    def initialize(
      @id : Int64?,
      @name : String?
    )
    end
  end
end
