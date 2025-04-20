module Audiofeeler
  struct Deploy
    getter id, account_id, server, local_dir, remote_dir, username, password

    def initialize(
      @id : Int64? = nil,
      @account_id : Int64? = nil,
      @server : String? = nil,
      @local_dir : String? = nil,
      @remote_dir : String? = nil,
      @username : String? = nil,
      @password : String? = nil,
    )
    end
  end
end
