package main

import (
    "context"
    "fmt"
    "os"
    "github.com/jackc/pgx/v5"
    "time"
    "github.com/mradmacher/audiofeeler/internal/repo"
    "github.com/mradmacher/audiofeeler/optiomist"
)

type Event struct {
    ID      uint32
    Date    string
    Hour    string
    Venue   string
    Address string
    Town    string
}

func main() {
    db, err := pgx.Connect(
        context.Background(),
        os.Getenv("AUDIOFEELER_DATABASE_URL"),
    )
    if err != nil {
        panic(err)
    }
    defer db.Close(context.Background())

    r := repo.EventsRepo { db }

    if r.Ping() {
        fmt.Println("Connected")
    } else {
        fmt.Println("Not connected")
    }

    _, err = db.Exec(context.Background(),
        `
        CREATE TABLE IF NOT EXISTS events (
            id SERIAL PRIMARY KEY,
            date date,
            hour time,
            venue VARCHAR(255),
            address VARCHAR(255),
            town VARCHAR(255)
        );
        `,
    )

    if err != nil {
        panic(err)
    }

    defer func() {
        _, err = db.Exec(context.Background(), "DROP TABLE IF EXISTS events;")
        fmt.Println("Tables dropped")
    }()

    fmt.Println("Tables created")

    params := repo.EventParams {
        Date: optiomist.Some(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
        Hour: optiomist.Nil[time.Time](),
        Venue: optiomist.Some("Some venue"),
        Address: optiomist.Some("Some address"),
        Town: optiomist.Some("Some town"),

    }
    id, err := r.Create(params)
    fmt.Printf("Event created [%v]: %v\n", id, err)

    params = repo.EventParams {
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
