package audiofeeler

import (
	. "github.com/mradmacher/audiofeeler/pkg/optiomist"
	"gopkg.in/yaml.v3"
	"io"
)

type jsonEvent struct {
	Date    string
	Hour    string
	Venue   string
	Address string
	City    string
	Place   string
}

func LoadEvents(reader io.Reader) (events []EventParams, err error) {
	var jsonEvents []jsonEvent

	decoder := yaml.NewDecoder(reader)
	err = decoder.Decode(&jsonEvents)
	if err != nil {
		return nil, err
	}
	for _, jEvent := range jsonEvents {
		var event = EventParams{}
		if jEvent.Date != "" {
			event.Date = Some(jEvent.Date)
		} else {
			event.Date = None[string]()
		}
		if jEvent.Hour != "" {
			event.Hour = Some(jEvent.Hour)
		} else {
			event.Hour = None[string]()
		}
		if jEvent.Venue != "" {
			event.Venue = Some(jEvent.Venue)
		} else {
			event.Venue = None[string]()
		}
		if jEvent.Address != "" {
			event.Address = Some(jEvent.Address)
		} else {
			event.Address = None[string]()
		}
		if jEvent.City != "" {
			event.City = Some(jEvent.City)
		} else {
			event.City = None[string]()
		}
		if jEvent.Place != "" {
			event.City = Some(jEvent.Place)
		} else {
			event.Place = None[string]()
		}
		events = append(events, event)
	}
	return events, nil
}
