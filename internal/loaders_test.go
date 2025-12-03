package audiofeeler

import (
	"github.com/mradmacher/audiofeeler/pkg/optiomist"
	"strings"
	"testing"
	"time"
)

var exampleEvents string = `
    - date: 2023-11-24
      hour: 20:00
      venue: Klub XYZ
      address: Mostowa 2
      town: Kraków
    - date: 2023-08-10
      hour: 19:30
      venue: Księgarnia podróżnicza ABC
      town: Kraków
    - date: 2024-01-01
      venue: Podgórska Jesień
`

func TestLoadEvents(t *testing.T) {
	reader := strings.NewReader(exampleEvents)
	events, err := LoadEvents(reader)

	if err != nil {
		t.Fatalf("Error while collecting events: %v", err)
	}

	if len(events) != 3 {
		t.Fatalf("Collected %d events; expected 3", len(events))
	}

	wants := []Event{
		Event{
			Date:    optiomist.Some(time.Date(2023, 11, 24, 0, 0, 0, 0, time.UTC)),
			Hour:    optiomist.Some(time.Date(0, 1, 1, 20, 0, 0, 0, time.UTC)),
			Venue:   optiomist.Some("Klub XYZ"),
			Address: optiomist.Some("Mostowa 2"),
			Town:    optiomist.Some("Kraków"),
		},
		Event{
			Date:    optiomist.Some(time.Date(2023, 8, 10, 0, 0, 0, 0, time.UTC)),
			Hour:    optiomist.Some(time.Date(0, 1, 1, 19, 30, 0, 0, time.UTC)),
			Venue:   optiomist.Some("Księgarnia podróżnicza ABC"),
			Address: optiomist.None[string](),
			Town:    optiomist.Some("Kraków"),
		},
		Event{
			Date:    optiomist.Some(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			Hour:    optiomist.None[time.Time](),
			Venue:   optiomist.Some("Podgórska Jesień"),
			Address: optiomist.None[string](),
			Town:    optiomist.None[string](),
		},
	}

	for i, want := range wants {
		if got := events[i]; got != want {
			t.Errorf("events[%d] = %+v; want %+v", i, got, want)
		}
	}
}
