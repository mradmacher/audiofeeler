package repo

import (
    "context"
    "strings"
    "time"
    "fmt"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgtype"
    "github.com/mradmacher/audiofeeler/optiomist"
)

type EventParams struct {
    ID      optiomist.Option[uint32]
    Date    optiomist.Option[time.Time]
    Hour    optiomist.Option[time.Time]
    Venue   optiomist.Option[string]
    Address optiomist.Option[string]
    Town    optiomist.Option[string]
}

type DbEvent struct {
    ID pgtype.Uint32
    Date pgtype.Date
    Hour pgtype.Time
    Venue pgtype.Text
    Address pgtype.Text
    Town pgtype.Text
}

type EventsRepo struct {
    Db *pgx.Conn
}

type Fields map[string]optiomist.Optionable

func buildInsert(tableName string, fields Fields) (string, []any) {
    //names := make([]string, 5)
    //params := make([]string, 5)
    //values := make([]any, 5)
    var names []string
    var params []string
    var values []any
    n := 1
    for name, value := range fields {
        if value.IsSome() {
            names = append(names, name)
            if value.IsNil() {
                params = append(params, "NULL")
            } else {
                params = append(params, fmt.Sprintf("$%d", n))
                n++
                values = append(values, value.Value())
            }
        }
    }
    query := "INSERT INTO " + tableName + "(" +
        strings.Join(names, ", ") + ") VALUES (" +
        strings.Join(params, ", ") + ") RETURNING id"
    return query, values
}

func (repo *EventsRepo) Create(params EventParams) (uint, error) {
    query, values := buildInsert(
        "events",
        Fields{
            "date": params.Date,
            "hour": params.Hour,
            "venue": params.Venue,
            "address": params.Address,
            "town": params.Town,
        },
    )
    row := repo.Db.QueryRow(context.Background(),
        query,
        values...,
    )
    var id uint
    err := row.Scan(&id)
    return id, err
}

func (repo *EventsRepo) All() (*[]EventParams, error) {
    var events []EventParams

    rows, err := repo.Db.Query(context.Background(),
        `
        SELECT id, date, hour, venue, address, town FROM events;
        `,
    )
    defer rows.Close()
    if err != nil {
        return nil, err
    }
    for rows.Next() {
        var event DbEvent
        err = rows.Scan(
            &event.ID,
            &event.Date,
            &event.Hour,
            &event.Venue,
            &event.Address,
            &event.Town,
        )

        eventParams := EventParams {
            ID: optiomist.Optiomize(event.ID.Uint32, event.ID.Valid),
            Date: optiomist.Optiomize(event.Date.Time, event.Date.Valid),
            Venue: optiomist.Optiomize(event.Venue.String, event.Venue.Valid),
            Address: optiomist.Optiomize(event.Address.String, event.Address.Valid),
            Town: optiomist.Optiomize(event.Town.String, event.Town.Valid),
        }
        opt := optiomist.Optiomize(event.Hour.Microseconds, event.Hour.Valid)
        if opt.IsSome() {
            eventParams.Hour = optiomist.Some(time.UnixMicro(opt.TypedValue()))
        } else {
            eventParams.Hour = optiomist.None[time.Time]()
        }

        events = append(events, eventParams)
    }

    if rows.Err() != nil {
        return nil, rows.Err()
    }

    return &events, nil
}

func (repo *EventsRepo) Ping() (bool) {

    err := repo.Db.Ping(context.Background())
    if err != nil {
        return false
    } else {
        return true
    }
}
