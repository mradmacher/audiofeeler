package repo

import (
    "context"
    "time"
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
    Db *DbClient
}

func (repo *EventsRepo) Create(params EventParams) (uint32, error) {
    fields := Fields{
		"date": params.Date,
		"hour": params.Hour,
		"venue": params.Venue,
		"address": params.Address,
		"town": params.Town,
	}
    query, values := fields.BuildInsert("events")
    row := repo.Db.Conn.QueryRow(context.Background(),
        query,
        values...,
    )
    var id uint32
    err := row.Scan(&id)
    return id, err
}

func buildEventParams(event DbEvent) *EventParams {
	eventParams := EventParams {
		ID: optiomist.Optiomize(event.ID.Uint32, event.ID.Valid),
		Date: optiomist.Optiomize(event.Date.Time, event.Date.Valid),
		Venue: optiomist.Optiomize(event.Venue.String, event.Venue.Valid),
		Address: optiomist.Optiomize(event.Address.String, event.Address.Valid),
		Town: optiomist.Optiomize(event.Town.String, event.Town.Valid),
	}
	opt := optiomist.Optiomize(event.Hour.Microseconds, event.Hour.Valid)
	if opt.IsSome() {
		eventParams.Hour = optiomist.Some(time.UnixMicro(opt.Value()))
	} else {
		eventParams.Hour = optiomist.None[time.Time]()
	}
	return &eventParams
}

func (repo *EventsRepo) Find(id uint32) (*EventParams, error) {
    row := repo.Db.Conn.QueryRow(context.Background(),
        `
        SELECT id, date, hour, venue, address, town
		FROM events
		WHERE id = $1
        `,
		id,
    )

	var event DbEvent
	err := row.Scan(
		&event.ID,
		&event.Date,
		&event.Hour,
		&event.Venue,
		&event.Address,
		&event.Town,
	)

    if err != nil {
        return nil, err
    }

	return buildEventParams(event), nil
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

        events = append(events, *buildEventParams(event))
    }

    if rows.Err() != nil {
        return nil, rows.Err()
    }

    return &events, nil
}

