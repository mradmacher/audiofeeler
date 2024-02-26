package repo

import (
    "context"
    "strings"
    "time"
    "fmt"
    "github.com/jackc/pgx/v5/pgxpool"
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

type DbClient struct {
    Conn *pgxpool.Pool
}

func Connect(dbUrl string) (*DbClient, error) {
    conn, err := pgxpool.New(context.Background(), dbUrl)
    if err != nil {
        return nil, err
    }
	return &DbClient{ conn }, nil
}

func (db *DbClient) Close() {
    db.Conn.Close()
}

func (db *DbClient) Ping() (bool) {
    err := db.Conn.Ping(context.Background())
    if err != nil {
        return false
    } else {
        return true
    }
}

type EventsRepo struct {
    Db *DbClient
}

type Fields map[string]optiomist.Optionable

func (db *DbClient) CreateStructure() error {
    _, err := db.Conn.Exec(
		context.Background(),
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
        return err
    }
	return nil
}

func (db *DbClient) RemoveStructure() error {
	_, err := db.Conn.Exec(
		context.Background(),
		`
		DROP TABLE IF EXISTS events;
		`,
	)
	return err
}

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
    query := "INSERT INTO " + tableName
	if len(names) > 0 {
		query = query + " (" +
        strings.Join(names, ", ") + ") VALUES (" +
        strings.Join(params, ", ") + ")"
	} else {
		query = query + " DEFAULT VALUES"
	}
	query = query + " RETURNING id"
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
    row := repo.Db.Conn.QueryRow(context.Background(),
        query,
        values...,
    )
    var id uint
    err := row.Scan(&id)
    return id, err
}

func (repo *EventsRepo) All() (*[]EventParams, error) {
    var events []EventParams

    rows, err := repo.Db.Conn.Query(context.Background(),
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

