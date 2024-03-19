package audiofeeler

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RecordNotFound struct{}

func (err RecordNotFound) Error() string {
	return "No record found in database"
}

func newRecordNotFound() error {
	return RecordNotFound{}
}

func wrapRecordNotFound(err error) error {
	if err == pgx.ErrNoRows {
		return newRecordNotFound()
	}
	return err
}

type DbClient struct {
	Conn *pgxpool.Pool
}

func NewDbClient(dbUrl string) (*DbClient, error) {
	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}
	return &DbClient{conn}, nil
}

func (db *DbClient) Close() {
	db.Conn.Close()
}

func (db *DbClient) Ping() bool {
	err := db.Conn.Ping(context.Background())
	if err != nil {
		return false
	} else {
		return true
	}
}

func (db *DbClient) CreateStructure() error {
	_, err := db.Conn.Exec(
		context.Background(),
		`
		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			title VARCHAR(255) NOT NULL,
			url VARCHAR(255) UNIQUE NOT NULL
		);

        CREATE TABLE IF NOT EXISTS events (
            id SERIAL PRIMARY KEY,
			account_id INTEGER NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
            date date,
            hour time,
            venue VARCHAR(255),
            address VARCHAR(255),
            town VARCHAR(255)
        );

		CREATE INDEX events_account_id_idx ON events (account_id);
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
		DROP INDEX IF EXISTS events_account_id_idx;
		DROP TABLE IF EXISTS events;
		DROP TABLE IF EXISTS accounts;
		`,
	)
	return err
}
