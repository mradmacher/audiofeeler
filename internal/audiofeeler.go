package audiofeeler

import (
	"github.com/mradmacher/audiofeeler/optiomist"
	"time"
)

type Event struct {
	ID      optiomist.Option[uint32]
	Date    optiomist.Option[time.Time]
	Hour    optiomist.Option[time.Time]
	Venue   optiomist.Option[string]
	Address optiomist.Option[string]
	Town    optiomist.Option[string]
}

type Artist struct {
	Name string
}

type Video struct {
	Url string
}
