package repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mradmacher/audiofeeler/internal"
	"github.com/mradmacher/audiofeeler/optiomist"
	"time"
)

type eventRecord struct {
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

func buildEventParams(record eventRecord) *audiofeeler.Event {
	event := audiofeeler.Event{
		Id:        optiomist.Optiomize(record.id.Uint32, record.id.Valid),
		AccountId: optiomist.Optiomize(record.accountId.Uint32, record.accountId.Valid),
		Date:      optiomist.Optiomize(record.date.Time, record.date.Valid),
		Venue:     optiomist.Optiomize(record.venue.String, record.venue.Valid),
		Address:   optiomist.Optiomize(record.address.String, record.address.Valid),
		Town:      optiomist.Optiomize(record.town.String, record.town.Valid),
	}
	opt := optiomist.Optiomize(record.hour.Microseconds, record.hour.Valid)
	if opt.IsSome() {
		event.Hour = optiomist.Some(time.UnixMicro(opt.Value()))
	} else {
		event.Hour = optiomist.None[time.Time]()
	}
	return &event
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

	var record eventRecord
	err := row.Scan(
		&record.id,
		&record.accountId,
		&record.date,
		&record.hour,
		&record.venue,
		&record.address,
		&record.town,
	)

	if err != nil {
		return nil, err
	}

	return buildEventParams(record), nil
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
		var record eventRecord
		err = rows.Scan(
			&record.id,
			&record.accountId,
			&record.date,
			&record.hour,
			&record.venue,
			&record.address,
			&record.town,
		)

		events = append(events, *buildEventParams(record))
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return &events, nil
}
