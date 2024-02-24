package audiofeeler

type Event struct {
	Date    string
	Hour    string
	Url     string
	Venue   string
	Address string
	Town    string
}

type Artist struct {
	Name string
}

type Video struct {
	Url string
}
