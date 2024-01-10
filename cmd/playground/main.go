package main

import (
  "gopkg.in/yaml.v3"
  "fmt"
  "github.com/mradmacher/audiofeeler/internal"
)

func loadData(filePath string) ([]audiofeeler.Event, error) {
	var events []audiofeeler.Event
	yamlBlob, err := os.ReadFile(filePath)

	err = yaml.Unmarshal(jsonBlob, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func seedDb() {

	var events []audiofeeler.Event
	events, err = loadData("events.yaml")
	if err != nil {
		panic(err)
	}

	for _, event := range events {
    fmt.Printf("%v\n", event)
	}
}

func main() {
  artist := audiofeeler.Artist{
    "Czarny motyl",
  }
  event := audiofeeler.Event{
    Date: "24.10.2023",
    Url: "http://example.com",
    Locations: []string{"Place", "City", "Country"},
  }

  fmt.Printf("%v\n", artist)
  fmt.Printf("%v\n", event)
}
