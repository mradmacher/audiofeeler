package audiofeeler

import (
	"os"
	"testing"
)

func setupTest(t *testing.T) (func(*testing.T), *DbClient) {
	db, err := NewDbClient(os.Getenv("AUDIOFEELER_TEST_DATABASE_URL"))
	if err != nil {
		t.Fatalf("Can't connect to DB: %v", err)
	}
	err = db.CreateStructure()
	if err != nil {
		t.Fatalf("Can't create tables: %v", err)
	}
	return func(t *testing.T) {
		db.RemoveStructure()
		db.Close()
	}, db
}
