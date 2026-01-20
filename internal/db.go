package audiofeeler

import (
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseId string
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

func NewDatabaseId() (DatabaseId, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return NewUnsetDatabaseId(), err
	}

	return DatabaseId(id.String()), nil
}

func NewUnsetDatabaseId() DatabaseId {
	return DatabaseId("")
}

func IsDatabaseIdSet(id DatabaseId) bool {
	if id == NewUnsetDatabaseId() {
		return false
	} else {
		return true
	}
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

func (db *DbClient) CreateStructure() error {
	_, err := db.Conn.Exec(
		`
		CREATE TABLE account (
		  id TEXT PRIMARY KEY,
		  name TEXT
		);

		CREATE TABLE deployment (
		  id TEXT PRIMARY KEY,
		  account_id INTEGER NOT NULL,
		  server TEXT,
		  username TEXT,
		  username_iv TEXT,
		  password TEXT,
		  password_iv TEXT,
		  remote_dir TEXT,
		  FOREIGN KEY(account_id) REFERENCES account(id) ON UPDATE RESTRICT ON DELETE RESTRICT
		);

		CREATE TABLE event (
		  id TEXT PRIMARY KEY,
		  account_id INTEGER NOT NULL,
		  name TEXT,
		  date TEXT,
		  hour TEXT,
		  venue TEXT,
		  town TEXT,
		  location TEXT,
		  description TEXT,
		  status INTEGER,
		  FOREIGN KEY(account_id) REFERENCES account(id) ON UPDATE RESTRICT ON DELETE RESTRICT
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
		DROP TABLE IF EXISTS event;
		DROP TABLE IF EXISTS deployment;
		DROP TABLE IF EXISTS account;
		`,
	)
	return err
}
