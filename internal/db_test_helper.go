package internal

import (
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func setupDbTest(t *testing.T) (func(*testing.T), DbEngine) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Can't load .env file")
	}
	db, err := NewDbEngine(os.Getenv("AUDIOFEELER_TEST_DATABASE_URL"))
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
