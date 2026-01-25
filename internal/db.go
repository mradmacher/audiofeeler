package audiofeeler

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

const sqlDir = "../db/"

type DatabaseId string

var ErrNotFound = errors.New("Record with given ID not found in database")
var ErrIdGenFailed = errors.New("Database ID generation failed")

func errNotFoundOr(err error) error {
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	return err
}

type DbClient struct {
	Conn *sql.DB
}

func NewDatabaseId() (DatabaseId, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return NewUnsetDatabaseId(), errors.Join(ErrIdGenFailed, err)
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
	path := filepath.Join(sqlDir, "schema.sql")
	schema, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = db.Conn.Exec(string(schema))

	if err != nil {
		return err
	}
	return nil
}

func (db *DbClient) RemoveStructure() error {
	path := filepath.Join(sqlDir, "drop.sql")
	schema, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = db.Conn.Exec(string(schema))
	return err
}
