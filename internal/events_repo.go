package audiofeeler

type DatabaseId uint32
type EventStatus int

const (
	EventCurrent EventStatus = iota
	EventArchived
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
	id          uint32
	accountId   uint32
	name        string
	date        string
	hour        string
	venue       string
	town        string
	location    string
	description string
	status      int
}

type EventsRepo struct {
	Db *DbClient
}

func (repo *EventsRepo) Save(event Event) (DatabaseId, error) {
	var query string
	if event.Id == 0 {
		query = "INSERT INTO events (account_id, name, date, hour, venue, town, location, description, status) " +
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) " +
			"RETURNING id;"
	}
	row := repo.Db.Conn.QueryRow(
		query,
		event.AccountId, event.Name, event.Date, event.Hour, event.Venue, event.Town, event.Location, event.Description, event.Status,
	)
	var id uint32
	err := row.Scan(&id)
	return DatabaseId(id), err
}

func buildEvent(record eventRecord) *Event {
	event := Event{
		Id:          DatabaseId(record.id),
		AccountId:   DatabaseId(record.accountId),
		Name:        record.name,
		Date:        record.date,
		Hour:        record.hour,
		Venue:       record.venue,
		Town:        record.town,
		Location:    record.location,
		Description: record.description,
		Status:      EventStatus(record.status),
	}
	return &event
}

func (repo *EventsRepo) Find(id DatabaseId) (*Event, error) {
	row := repo.Db.Conn.QueryRow(
		`
        SELECT id, account_id, name, date, hour, venue, town, location, description, status
		FROM events
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

	if err != nil {
		return nil, err
	}

	return buildEvent(record), nil
}

func (repo *EventsRepo) FindAll(accountId DatabaseId) ([]Event, error) {
	var events []Event
	rows, err := repo.Db.Conn.Query(
		`
        SELECT id, account_id, name, date, hour, venue, town, location, description, status
		FROM events
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


		events = append(events, *buildEvent(record))
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return events, nil
}
