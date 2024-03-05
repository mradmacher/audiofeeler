package repo

import (
	"os"
	"testing"
)

func setupTest(t *testing.T) (func(*testing.T), *DbClient) {
	db, err := Connect(os.Getenv("AUDIOFEELER_TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal("Can't connect to DB")
	}
	err = db.CreateStructure()
	if err != nil {
		t.Fatal("Can't create tables")
	}
	return func(t *testing.T) {
		db.RemoveStructure()
		db.Close()
	}, db
}
