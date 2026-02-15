package internal

type EventStatus int

const (
	CurrentEvent EventStatus = iota
	ArchivedEvent
)

type Event struct {
	Id          DatabaseId
	AccountId   DatabaseId
	Name        string
	Date        string
	Hour        string
	Venue       string
	Town        string
	Location    string
	Description string
	Status      EventStatus
}

type eventRecord struct {
	id          string
	accountId   string
	name        string
	date        string
	hour        string
	venue       string
	town        string
	location    string
	description string
	status      string
}

type EventRepo struct {
	Db DbEngine
}

func (repo *EventRepo) Save(event Event) (DatabaseId, error) {
	var query string
	var id DatabaseId
	var err error
	var eventStatus string
	if event.Status == CurrentEvent {
		eventStatus = "current"
	} else {
		eventStatus = "archived"
	}
	if IsDatabaseIdSet(event.Id) {
		id = event.Id
		query = "UPDATE event SET name = $1, date = $2, hour = $3, venue = $4, town = $5, location = $6, description = $7, status = $8 " +
			"WHERE id = $9;"
		_, err = repo.Db.Conn.Exec(
			query,
			event.Name, event.Date, event.Hour, event.Venue, event.Town, event.Location, event.Description, eventStatus, id,
		)
	} else {
		query = "INSERT INTO event (id, account_id, name, date, hour, venue, town, location, description, status) " +
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);"
		if id, err = NewDatabaseId(); err != nil {
			return id, err
		}
		_, err = repo.Db.Conn.Exec(
			query,
			id, event.AccountId, event.Name, event.Date, event.Hour, event.Venue, event.Town, event.Location, event.Description, eventStatus,
		)
	}
	if err != nil {
		return NewUnsetDatabaseId(), err
	}
	return id, err
}

func (repo *EventRepo) Delete(id DatabaseId) error {
	_, err := repo.Db.Conn.Exec("DELETE FROM event where ID = $1", id)
	return err
}

func buildEvent(record eventRecord, isFound bool) *Event {
	var event Event

	if isFound {
		event = Event{
			Id:          DatabaseId(record.id),
			AccountId:   DatabaseId(record.accountId),
			Name:        record.name,
			Date:        record.date,
			Hour:        record.hour,
			Venue:       record.venue,
			Town:        record.town,
			Location:    record.location,
			Description: record.description,
		}
		switch record.status {
		case "archived":
			event.Status = ArchivedEvent
		default:
			event.Status = CurrentEvent
		}
	} else {
		event = Event{}
	}

	return &event
}

func (repo *EventRepo) Find(id DatabaseId) (FindResult[Event], error) {
	row := repo.Db.Conn.QueryRow(
		`
        SELECT id, account_id, name, date, hour, venue, town, location, description, status
		FROM event
		WHERE id = $1
        `,
		id,
	)

	var record eventRecord
	err := row.Scan(
		&record.id,
		&record.accountId,
		&record.name,
		&record.date,
		&record.hour,
		&record.venue,
		&record.town,
		&record.location,
		&record.description,
		&record.status,
	)

	found, err := FilterNotFoundErr(err)

	return FindResult[Event]{*buildEvent(record, found), found}, err
}

func (repo *EventRepo) FindAll(accountId DatabaseId) ([]Event, error) {
	var events []Event
	rows, err := repo.Db.Conn.Query(
		`
        SELECT id, account_id, name, date, hour, venue, town, location, description, status
		FROM event
		WHERE account_id = $1
        `,
		accountId,
	)

	defer rows.Close()

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var record eventRecord
		rows.Scan(
			&record.id,
			&record.accountId,
			&record.name,
			&record.date,
			&record.hour,
			&record.venue,
			&record.town,
			&record.location,
			&record.description,
			&record.status,
		)

		events = append(events, *buildEvent(record, true))
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return events, nil
}
