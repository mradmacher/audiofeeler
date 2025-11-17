package audiofeeler

import (
	"github.com/mradmacher/audiofeeler/optiomist"
	"gopkg.in/yaml.v3"
	"io"
	"time"
)

type jsonEvent struct {
	Date    string
	Hour    string
	Venue   string
	Address string
	Town    string
}

func LoadEvents(reader io.Reader) (events []Event, err error) {
	var jsonEvents []jsonEvent

	decoder := yaml.NewDecoder(reader)
	err = decoder.Decode(&jsonEvents)
	if err != nil {
		return nil, err
	}
	for _, jEvent := range jsonEvents {
		var event = Event{}
		if jEvent.Date != "" {
			date, err := time.Parse(time.DateOnly, jEvent.Date)
			if err != nil {
				event.Date = optiomist.None[time.Time]()
			} else {
				event.Date = optiomist.Some(date)
			}
		} else {
			event.Date = optiomist.None[time.Time]()
		}
		if jEvent.Hour != "" {
			hour, err := time.Parse(time.TimeOnly, jEvent.Hour+":00")
			if err != nil {
				event.Hour = optiomist.None[time.Time]()
			} else {
				event.Hour = optiomist.Some(hour)
			}
		} else {
			event.Hour = optiomist.None[time.Time]()
		}
		if jEvent.Venue != "" {
			event.Venue = optiomist.Some(jEvent.Venue)
		} else {
			event.Venue = optiomist.None[string]()
		}
		if jEvent.Address != "" {
			event.Address = optiomist.Some(jEvent.Address)
		} else {
			event.Address = optiomist.None[string]()
		}
		if jEvent.Town != "" {
			event.Town = optiomist.Some(jEvent.Town)
		} else {
			event.Town = optiomist.None[string]()
		}
		events = append(events, event)
	}
	return events, nil
}
