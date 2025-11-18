package audiofeeler

import (
	"github.com/mradmacher/audiofeeler/optiomist"
	"github.com/mradmacher/audiofeeler/sqlbuilder"
	"time"
)

type eventRecord struct {
	id        uint32
	accountId uint32
	date      time.Time
	hour      time.Time
	venue     string
	address   string
	town      string
}

type EventsRepo struct {
	Db *DbClient
}

func (repo *EventsRepo) Create(event Event) (uint32, error) {
	fields := sqlbuilder.Fields{
		"account_id": event.AccountId,
		"date":       event.Date,
		"hour":       event.Hour,
		"venue":      event.Venue,
		"address":    event.Address,
		"town":       event.Town,
	}
	query, values := fields.BuildInsert("events")
	row := repo.Db.Conn.QueryRow(
		query,
		values...,
	)
	var id uint32
	err := row.Scan(&id)
	return id, err
}

func buildEventParams(record eventRecord) *Event {
	event := Event{
		Id:        optiomist.Optiomize(record.id, true),
		AccountId: optiomist.Optiomize(record.accountId, true),
		Date:      optiomist.Optiomize(record.date, true),
		Venue:     optiomist.Optiomize(record.venue, true),
		Address:   optiomist.Optiomize(record.address, true),
		Town:      optiomist.Optiomize(record.town, true),
		Hour:      optiomist.Optiomize(record.hour, true),
	}
	/*
	opt := optiomist.Optiomize(record.hour.Microseconds, record.hour.Valid)
	if opt.IsSome() {
		event.Hour = optiomist.Some(time.UnixMicro(opt.Value()))
	} else {
		event.Hour = optiomist.None[time.Time]()
	}
	*/
	return &event
}

func (repo *EventsRepo) Find(id uint32) (*Event, error) {
	row := repo.Db.Conn.QueryRow(
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

func (repo *EventsRepo) FindAll() (*[]Event, error) {
	var events []Event

	rows, err := repo.Db.Conn.Query(
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
