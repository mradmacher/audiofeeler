package repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mradmacher/audiofeeler/internal"
	"github.com/mradmacher/audiofeeler/optiomist"
	"time"
)

type dbEvent struct {
	ID      pgtype.Uint32
	Date    pgtype.Date
	Hour    pgtype.Time
	Venue   pgtype.Text
	Address pgtype.Text
	Town    pgtype.Text
}

type EventsRepo struct {
	Db *DbClient
}

func (repo *EventsRepo) Create(event audiofeeler.Event) (uint32, error) {
	fields := Fields{
		"date":    event.Date,
		"hour":    event.Hour,
		"venue":   event.Venue,
		"address": event.Address,
		"town":    event.Town,
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

func buildEventParams(event dbEvent) *audiofeeler.Event {
	aEvent := audiofeeler.Event{
		ID:      optiomist.Optiomize(event.ID.Uint32, event.ID.Valid),
		Date:    optiomist.Optiomize(event.Date.Time, event.Date.Valid),
		Venue:   optiomist.Optiomize(event.Venue.String, event.Venue.Valid),
		Address: optiomist.Optiomize(event.Address.String, event.Address.Valid),
		Town:    optiomist.Optiomize(event.Town.String, event.Town.Valid),
	}
	opt := optiomist.Optiomize(event.Hour.Microseconds, event.Hour.Valid)
	if opt.IsSome() {
		aEvent.Hour = optiomist.Some(time.UnixMicro(opt.Value()))
	} else {
		aEvent.Hour = optiomist.None[time.Time]()
	}
	return &aEvent
}

func (repo *EventsRepo) Find(id uint32) (*audiofeeler.Event, error) {
	row := repo.Db.Conn.QueryRow(context.Background(),
		`
        SELECT id, date, hour, venue, address, town
		FROM events
		WHERE id = $1
        `,
		id,
	)

	var event dbEvent
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

func (repo *EventsRepo) All() (*[]audiofeeler.Event, error) {
	var events []audiofeeler.Event

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
		var event dbEvent
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
