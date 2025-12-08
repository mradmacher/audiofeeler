package audiofeeler

type Account struct {
	Id    int64
	Name  string
	SourceDir string
}

type Deployment struct {
	Id int64
	AccountId int64
	Server string
	Username string
	UsernameIV string
	Password string
	PasswordIV string
	RemoteDir string
}

type Event struct {
	Id        uint32
	AccountId int64
	Date      string
	Hour      string
	Venue     string
	Place     string
	City      string
	Address   string
}
