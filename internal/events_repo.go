package audiofeeler

import (
	. "github.com/mradmacher/audiofeeler/pkg/optiomist"
	"github.com/mradmacher/audiofeeler/pkg/sqlbuilder"
)
type EventParams struct {
	Id        Option[uint32]
	AccountId Option[int64]
	Date      Option[string]
	Hour      Option[string]
	Venue     Option[string]
	Place     Option[string]
	City      Option[string]
	Address   Option[string]
}

type eventRecord struct {
	id        uint32
	accountId int64
	date      string
	hour      string
	venue     string
	place     string
	city      string
	address   string
}

type EventsRepo struct {
	Db *DbClient
}

func (repo *EventsRepo) Create(event EventParams) (uint32, error) {
	fields := sqlbuilder.Fields{
		"account_id": event.AccountId,
		"date":       event.Date,
		"hour":       event.Hour,
		"venue":      event.Venue,
		"address":    event.Address,
		"city":       event.City,
		"place":      event.Place,
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
		Id:        record.id,
		AccountId: record.accountId,
		Date:      record.date,
		Hour:      record.hour,
		Venue:     record.venue,
		Address:   record.address,
		City:      record.city,
		Place:     record.place,
	}
	return &event
}

func (repo *EventsRepo) Find(id uint32) (*Event, error) {
	row := repo.Db.Conn.QueryRow(
		`
        SELECT id, account_id, date, hour, venue, address, city, place
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
		&record.city,
		&record.place,
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
        SELECT id, account_id, date, hour, venue, address, city, place FROM events;
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
			&record.city,
			&record.place,
		)

		events = append(events, *buildEventParams(record))
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return &events, nil
}
