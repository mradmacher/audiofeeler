package audiofeeler
/*
import (
	. "github.com/mradmacher/audiofeeler/pkg/optiomist"
	"strings"
	"testing"
)

var exampleEvents string = `
    - date: 2023-11-24
      hour: 20:00
      name: Festiwal
      location: Mostowa 2
      town: Kraków
	  venue: Pub XYZ
    - date: 2023-08-10
      hour: 19:30
      venue: Księgarnia podróżnicza ABC
      town: Kraków
    - date: 2024-01-01
      name: Podgórska Jesień
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

	wants := []EventParams{
		EventParams{
			Date:    Some("2023-11-24"),
			Hour:    Some("20:00"),
			Venue:   Some("Festiwal"),
			Address: Some("Mostowa 2"),
			City:    Some("Kraków"),
			Place:   Some("Pub XYZ"),
		},
		EventParams{
			Date:    Some("2023-08-10"),
			Hour:    Some("19:30"),
			Venue:   Some("Księgarnia podróżnicza ABC"),
			Address: None[string](),
			City:    Some("Kraków"),
			Place:   None[string](),
		},
		EventParams{
			Date:    Some("2024-01-01"),
			Hour:    None[string](),
			Venue:   Some("Podgórska Jesień"),
			Address: None[string](),
			City:    None[string](),
			Place:   None[string](),
		},
	}

	for i, want := range wants {
		if got := events[i]; got != want {
			t.Errorf("events[%d] = %+v; want %+v", i, got, want)
		}
	}
}
*/
