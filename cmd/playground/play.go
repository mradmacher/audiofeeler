package main

import (
    "fmt"
    "os"
    "time"
    "github.com/mradmacher/audiofeeler/internal/repo"
    "github.com/mradmacher/audiofeeler/optiomist"
    "github.com/mradmacher/audiofeeler/internal"
)

func main() {
    db, err := repo.Connect(os.Getenv("AUDIOFEELER_DATABASE_URL"))
    if err != nil {
        panic(err)
    }
    defer db.Close()

    r := repo.EventsRepo { db }

    if db.Ping() {
        fmt.Println("Connected")
    } else {
        fmt.Println("Not connected")
    }

	err = db.CreateStructure()
    if err != nil {
        panic(err)
    }

    defer func() {
        db.RemoveStructure()
        fmt.Println("Tables dropped")
    }()

    fmt.Println("Tables created")

    params := audiofeeler.Event {
        Date: optiomist.Some(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
        Hour: optiomist.Nil[time.Time](),
        Venue: optiomist.Some("Some venue"),
        Address: optiomist.Some("Some address"),
        Town: optiomist.Some("Some town"),

    }
    id, err := r.Create(params)
    fmt.Printf("Event created [%v]: %v\n", id, err)

    params = audiofeeler.Event {
        Date: optiomist.Some(time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)),
        Hour: optiomist.Some(time.Date(0, 0, 0, 21, 0, 0, 0, time.UTC)),
        Venue: optiomist.Some("Other venue"),
        Address: optiomist.Some("Other address"),
        Town: optiomist.None[string](),

    }
    id, err = r.Create(params)
    fmt.Printf("Event created [%v]: %v\n", id, err)

    events, err := r.All()
    if err != nil {
        panic(err)
    }

    for _, event := range *events {
        fmt.Printf("Event: %v\n", event)
    }
}
