package audiofeeler

type Deployment struct {
	Id         int64
	AccountId  int64
	Server     string
	Username   string
	UsernameIV string
	Password   string
	PasswordIV string
	RemoteDir  string
}
