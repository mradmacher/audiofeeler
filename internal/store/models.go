package store

import (
	"github.com/mradmacher/audiofeeler/optiomist"
	"time"
)

type Account struct {
	Id    optiomist.Option[uint32]
	Name  optiomist.Option[string]
	Title optiomist.Option[string]
	Url   optiomist.Option[string]
}

type Event struct {
	Id        optiomist.Option[uint32]
	AccountId optiomist.Option[uint32]
	Date      optiomist.Option[time.Time]
	Hour      optiomist.Option[time.Time]
	Venue     optiomist.Option[string]
	Address   optiomist.Option[string]
	Town      optiomist.Option[string]
}

type Artist struct {
	Name string
}

type Video struct {
	Url string
}
