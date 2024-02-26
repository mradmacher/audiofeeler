package repo

import (
	"testing"
	"os"
    "time"
    "github.com/mradmacher/audiofeeler/optiomist"
)

func setup(t *testing.T) (func(*testing.T), *DbClient) {
	db, err := Connect(os.Getenv("AUDIOFEELER_TEST_DATABASE_URL"))
    if err != nil {
		t.Fatal("Can't connect to DB")
    }
	err = db.CreateStructure()
    if err != nil {
		t.Fatal("Can't create tables")
    }
	return func(t *testing.T) {
		//db.RemoveStructure()
		db.Close()
	}, db
}

func TestCreate(t *testing.T) {
	teardown, db := setup(t)
	defer teardown(t)

    r := EventsRepo { db }

	tests := []struct {
		name string
		params EventParams
	}{
		{
			"some params",

			EventParams {
				Date: optiomist.Some(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
				Hour: optiomist.Some(time.Date(0, 0, 0, 21, 0, 0, 0, time.UTC)),
				Venue: optiomist.Some("Some venue"),
				Address: optiomist.Some("Some address"),
				Town: optiomist.Some("Some town"),
			},
		}, {
			"none params",
			EventParams {
				Date: optiomist.None[time.Time](),
				Hour: optiomist.None[time.Time](),
				Venue: optiomist.None[string](),
				Address: optiomist.None[string](),
				Town: optiomist.None[string](),
			},
		}, {
			"nil params",
			EventParams {
				Date: optiomist.Nil[time.Time](),
				Hour: optiomist.Nil[time.Time](),
				Venue: optiomist.Nil[string](),
				Address: optiomist.Nil[string](),
				Town: optiomist.Nil[string](),
			},
		},
    }

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			id, err := r.Create(test.params)

			if err != nil {
				t.Fatal(err)
			}
			if id <= 0 {
				t.Error("Id is less than zero")
			}
		})
	}
}
