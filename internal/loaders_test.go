package audiofeeler

import (
	"strings"
	"testing"
)

var exampleEvents string = `
    - date: 24.11.2023
      hour: 20:00
      venue: Klub XYZ
      address: Mostowa 2
      town: Kraków
      url: https://www.example.com/events/xyz
    - date: 10.08.2023
      hour: 19:30
      venue: Księgarnia podróżnicza ABC
      town: Kraków
    - date: 01.01.2024
      venue: Podgórska Jesień
`

func TestLoadEvents_fetchesAllEvents(t *testing.T) {
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
			Date:    "24.11.2023",
			Hour:    "20:00",
			Venue:   "Klub XYZ",
			Address: "Mostowa 2",
			Town:    "Kraków",
			Url:     "https://www.example.com/events/xyz",
		},
		Event{
			Date:  "10.08.2023",
			Hour:  "19:30",
			Venue: "Księgarnia podróżnicza ABC",
			Town:  "Kraków",
		},
		Event{
			Date:  "01.01.2024",
			Venue: "Podgórska Jesień",
		},
	}

	for i, want := range wants {
		if got := events[i]; got != want {
			t.Errorf("events[%d] = %v; want %v", i, got, want)
		}
	}
}
