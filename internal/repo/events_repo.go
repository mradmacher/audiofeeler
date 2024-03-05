package repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mradmacher/audiofeeler/internal"
	"github.com/mradmacher/audiofeeler/optiomist"
	"time"
)

type dbEvent struct {
	id        pgtype.Uint32
	accountId pgtype.Uint32
	date      pgtype.Date
	hour      pgtype.Time
	venue     pgtype.Text
	address   pgtype.Text
	town      pgtype.Text
}

type EventsRepo struct {
	Db *DbClient
}

func (repo *EventsRepo) Create(event audiofeeler.Event) (uint32, error) {
	fields := Fields{
		"account_id": event.AccountId,
		"date":       event.Date,
		"hour":       event.Hour,
		"venue":      event.Venue,
		"address":    event.Address,
		"town":       event.Town,
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
		Id:        optiomist.Optiomize(event.id.Uint32, event.id.Valid),
		AccountId: optiomist.Optiomize(event.accountId.Uint32, event.accountId.Valid),
		Date:      optiomist.Optiomize(event.date.Time, event.date.Valid),
		Venue:     optiomist.Optiomize(event.venue.String, event.venue.Valid),
		Address:   optiomist.Optiomize(event.address.String, event.address.Valid),
		Town:      optiomist.Optiomize(event.town.String, event.town.Valid),
	}
	opt := optiomist.Optiomize(event.hour.Microseconds, event.hour.Valid)
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
        SELECT id, account_id, date, hour, venue, address, town
		FROM events
		WHERE id = $1
        `,
		id,
	)

	var event dbEvent
	err := row.Scan(
		&event.id,
		&event.accountId,
		&event.date,
		&event.hour,
		&event.venue,
		&event.address,
		&event.town,
	)

	if err != nil {
		return nil, err
	}

	return buildEventParams(event), nil
}

func (repo *EventsRepo) FindAll() (*[]audiofeeler.Event, error) {
	var events []audiofeeler.Event

	rows, err := repo.Db.Conn.Query(context.Background(),
		`
        SELECT id, account_id, date, hour, venue, address, town FROM events;
        `,
	)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var event dbEvent
		err = rows.Scan(
			&event.id,
			&event.accountId,
			&event.date,
			&event.hour,
			&event.venue,
			&event.address,
			&event.town,
		)

		events = append(events, *buildEventParams(event))
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &events, nil
}
