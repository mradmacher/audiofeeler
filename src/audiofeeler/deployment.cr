module Audiofeeler
  struct Deployment
    getter id, account_id, server, remote_dir, username, password

    def initialize(
      @id : Int64? = nil,
      @account_id : Int64? = nil,
      @server : String? = nil,
      @remote_dir : String? = nil,
      @username : String? = nil,
      @password : String? = nil,
    )
    end
  end
end
