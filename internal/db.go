package audiofeeler

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type RecordNotFound struct{}

func (err RecordNotFound) Error() string {
	return "No record found in database"
}

func newRecordNotFound() error {
	return RecordNotFound{}
}

func wrapRecordNotFound(err error) error {
	if err == sql.ErrNoRows {
		return newRecordNotFound()
	}
	return err
}

type DbClient struct {
	Conn *sql.DB
}

func NewDbClient(dbUrl string) (*DbClient, error) {
	conn, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		return nil, err
	}
	return &DbClient{conn}, nil
}

func (db *DbClient) Close() {
	db.Conn.Close()
}

/*
func (db *DbClient) Ping() bool {
	err := db.Conn.Ping(context.Background())
	if err != nil {
		return false
	} else {
		return true
	}
}
*/

func (db *DbClient) CreateStructure() error {
	_, err := db.Conn.Exec(
		`
		CREATE TABLE accounts (
		  id INTEGER PRIMARY KEY,
		  name TEXT
		);

		CREATE TABLE deployments (
		  id INTEGER PRIMARY KEY,
		  account_id INTEGER,
		  server TEXT,
		  username TEXT,
		  username_iv TEXT,
		  password TEXT,
		  password_iv TEXT,
		  remote_dir TEXT,
		  FOREIGN KEY(account_id) REFERENCES accounts(id)
		);

		CREATE TABLE events (
		  id INTEGER PRIMARY KEY,
		  account_id INTEGER,
		  name TEXT,
		  date TEXT,
		  hour TEXT,
		  venue TEXT,
		  town TEXT,
		  location TEXT,
		  description TEXT,
		  status INTEGER,
		  FOREIGN KEY(account_id) REFERENCES accounts(id)
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
		`
		DROP TABLE IF EXISTS events;
		DROP TABLE IF EXISTS deployments;
		DROP TABLE IF EXISTS accounts;
		`,
	)
	return err
}
