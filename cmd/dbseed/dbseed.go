package main

import (
	"bytes"
    "fmt"
    "os"
    "github.com/mradmacher/audiofeeler/internal/repo"
    "github.com/mradmacher/audiofeeler/internal"
)

func main() {
    if len(os.Args) < 2 {
        panic("Missing file name to parse")
    }
    fileName := os.Args[1]

    db, err := repo.Connect(os.Getenv("AUDIOFEELER_DATABASE_URL"))
    if err != nil {
        panic(err)
    }
    defer db.Close()
    if db.Ping() {
        fmt.Println("Connected to database")
    } else {
        panic("Not connected to database")
    }

	err = db.CreateStructure()
    if err != nil {
        panic(err)
    }
    fmt.Println("Tables created")

	jsonBlob, err := os.ReadFile(fileName)
	events, err := audiofeeler.LoadEvents(bytes.NewReader(jsonBlob))
    if err != nil {
        panic(err)
    }

    r := repo.EventsRepo { db }

	for _, event := range events {
		id, err := r.Create(event)
		fmt.Printf("Event created [%v]: %v\n", id, err)
	}
}
