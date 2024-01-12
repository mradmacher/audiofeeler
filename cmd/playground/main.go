package main

import (
	"fmt"
	"github.com/mradmacher/audiofeeler/internal"
	"gopkg.in/yaml.v3"
	"os"
)

func loadData(filePath string) (events []audiofeeler.Event, err error) {
	yamlBlob, err := os.ReadFile(filePath)

	err = yaml.Unmarshal(yamlBlob, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func seedDb() {
	var events []audiofeeler.Event
	var err error
	events, err = loadData("events.yml")
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
		Date:      "24.10.2023",
		Url:       "http://example.com",
		Locations: []string{"Place", "City", "Country"},
	}

	fmt.Printf("%v\n", artist)
	fmt.Printf("%v\n", event)

	seedDb()
}
