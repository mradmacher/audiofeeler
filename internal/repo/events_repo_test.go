package repo

import (
	"github.com/mradmacher/audiofeeler/internal"
	"github.com/mradmacher/audiofeeler/optiomist"
	"os"
	"testing"
	"time"
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
		db.RemoveStructure()
		db.Close()
	}, db
}

func TestEventRepo(t *testing.T) {
	teardown, db := setup(t)
	defer teardown(t)

	r := EventsRepo{db}

	t.Run("Create", testEventsRepo_Create(&r))
	t.Run("Find not nil values", testEventsRepo_Find_not_nils(&r))
	t.Run("Find nil values", testEventsRepo_Find_nils(&r))
}

func testEventsRepo_Create(r *EventsRepo) func(*testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name  string
			event audiofeeler.Event
		}{
			{
				"some params",
				audiofeeler.Event{
					Date:    optiomist.Some(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
					Hour:    optiomist.Some(time.Date(0, 0, 0, 21, 0, 0, 0, time.UTC)),
					Venue:   optiomist.Some("Some venue"),
					Address: optiomist.Some("Some address"),
					Town:    optiomist.Some("Some town"),
				},
			}, {
				"none params",
				audiofeeler.Event{
					Date:    optiomist.None[time.Time](),
					Hour:    optiomist.None[time.Time](),
					Venue:   optiomist.None[string](),
					Address: optiomist.None[string](),
					Town:    optiomist.None[string](),
				},
			}, {
				"nil params",
				audiofeeler.Event{
					Date:    optiomist.Nil[time.Time](),
					Hour:    optiomist.Nil[time.Time](),
					Venue:   optiomist.Nil[string](),
					Address: optiomist.Nil[string](),
					Town:    optiomist.Nil[string](),
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				id, err := r.Create(test.event)

				if err != nil {
					t.Fatal(err)
				}
				if id <= 0 {
					t.Error("Id needs to be a positive number")
				}
			})
		}
	}
}

func testEventsRepo_Find_not_nils(r *EventsRepo) func(*testing.T) {
	return func(t *testing.T) {
		event := audiofeeler.Event{
			Date:    optiomist.Some(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
			Hour:    optiomist.Some(time.Date(0, 0, 0, 21, 0, 0, 0, time.UTC)),
			Venue:   optiomist.Some("Some venue"),
			Address: optiomist.Some("Some address"),
			Town:    optiomist.Some("Some town"),
		}

		id, err := r.Create(event)
		if err != nil {
			t.Fatal(err)
		}
		params, err := r.Find(id)
		if err != nil {
			t.Fatal(err)
		}
		if params.ID.Value() != id {
			t.Errorf("ID = %v; expected %v", params.ID.Value(), id)
		}
		if params.Date != event.Date {
			t.Errorf("Date = %v; expected %v", params.Date.Value(), event.Date.Value())
		}
		if !params.Hour.IsSome() || !event.Hour.IsSome() || params.Hour.Value().Format(time.TimeOnly) != event.Hour.Value().Format(time.TimeOnly) {
			t.Errorf("Hour = %v; expected %v", params.Hour.Value(), event.Hour.Value())
		}
		if params.Venue != event.Venue {
			t.Errorf("Venue = %v; expected %v", params.Venue.Value(), event.Venue.Value())
		}
		if params.Address != event.Address {
			t.Errorf("Address = %v; expected %v", params.Address.Value(), event.Address.Value())
		}
		if params.Town != event.Town {
			t.Errorf("Town = %v; expected %v", params.Town.Value(), event.Town.Value())
		}
	}
}

func testEventsRepo_Find_nils(r *EventsRepo) func(*testing.T) {
	return func(t *testing.T) {
		event := audiofeeler.Event{
			Date:    optiomist.None[time.Time](),
			Hour:    optiomist.Nil[time.Time](),
			Venue:   optiomist.None[string](),
			Address: optiomist.Nil[string](),
			Town:    optiomist.None[string](),
		}

		id, err := r.Create(event)
		if err != nil {
			t.Fatal(err)
		}
		found, err := r.Find(id)
		if err != nil {
			t.Fatal(err)
		}
		if found.ID.Value() != id {
			t.Errorf("ID = %v; expected %v", found.ID.Value(), id)
		}
		if found.Date.IsSome() {
			t.Errorf("Date = %v; expected none", found.Date.Value())
		}
		if found.Hour.IsSome() {
			t.Errorf("Hour = %v; expected none", found.Hour.Value())
		}
		if found.Venue.IsSome() {
			t.Errorf("Venue = %v; expected none", found.Venue.Value())
		}
		if found.Address.IsSome() {
			t.Errorf("Address = %v; expected none", found.Address.Value())
		}
		if found.Town.IsSome() {
			t.Errorf("Town = %v; expected none", found.Town.Value())
		}
	}
}
