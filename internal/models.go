package internal

type FindResult[T any] struct {
	Record  T
	IsFound bool
}

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
