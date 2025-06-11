module Audiofeeler
  struct Account
    getter id, name, source_dir

    def initialize(
      @id : Int64?,
      @name : String?,
      @source_dir : String?
    )
    end
  end
end
