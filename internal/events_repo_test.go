package audiofeeler

import (
	"github.com/mradmacher/audiofeeler/optiomist"
	"testing"
	"time"
)

func setupAccount(db *DbClient, t *testing.T) uint32 {
	accountsRepo := AccountsRepo{db}
	accountId, err := accountsRepo.Create(Account{
		Name:  optiomist.Some("example"),
		Title: optiomist.Some("Example"),
		Url:   optiomist.Some("http://example.com"),
	})
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	return accountId
}

func TestEventsRepo(t *testing.T) {
	teardown, db := setupTest(t)
	defer teardown(t)

	accountId := setupAccount(db, t)

	r := EventsRepo{db}

	t.Run("Create", testEventsRepo_Create(&r, accountId))
	t.Run("Find not nil values", testEventsRepo_Find_not_nils(&r, accountId))
	t.Run("Find nil values", testEventsRepo_Find_nils(&r, accountId))
}

func testEventsRepo_Create(r *EventsRepo, accountId uint32) func(*testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name  string
			event Event
		}{
			{
				"some params",
				Event{
					AccountId: optiomist.Some(accountId),
					Date:      optiomist.Some(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
					Hour:      optiomist.Some(time.Date(0, 0, 0, 21, 0, 0, 0, time.UTC)),
					Venue:     optiomist.Some("Some venue"),
					Address:   optiomist.Some("Some address"),
					Town:      optiomist.Some("Some town"),
				},
			}, {
				"none params",
				Event{
					AccountId: optiomist.Some(accountId),
					Date:      optiomist.None[time.Time](),
					Hour:      optiomist.None[time.Time](),
					Venue:     optiomist.None[string](),
					Address:   optiomist.None[string](),
					Town:      optiomist.None[string](),
				},
			}, {
				"nil params",
				Event{
					AccountId: optiomist.Some(accountId),
					Date:      optiomist.Nil[time.Time](),
					Hour:      optiomist.Nil[time.Time](),
					Venue:     optiomist.Nil[string](),
					Address:   optiomist.Nil[string](),
					Town:      optiomist.Nil[string](),
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

func testEventsRepo_Find_not_nils(r *EventsRepo, accountId uint32) func(*testing.T) {
	return func(t *testing.T) {
		event := Event{
			AccountId: optiomist.Some(accountId),
			Date:      optiomist.Some(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
			Hour:      optiomist.Some(time.Date(0, 0, 0, 21, 0, 0, 0, time.UTC)),
			Venue:     optiomist.Some("Some venue"),
			Address:   optiomist.Some("Some address"),
			Town:      optiomist.Some("Some town"),
		}

		id, err := r.Create(event)
		if err != nil {
			t.Fatal(err)
		}
		params, err := r.Find(id)
		if err != nil {
			t.Fatal(err)
		}
		if params.Id.Value() != id {
			t.Errorf("Id = %v; expected %v", params.Id.Value(), id)
		}
		if params.AccountId.Value() != accountId {
			t.Errorf("AccountId = %v; expected %v", params.AccountId.Value(), accountId)
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

func testEventsRepo_Find_nils(r *EventsRepo, accountId uint32) func(*testing.T) {
	return func(t *testing.T) {
		event := Event{
			AccountId: optiomist.Some(accountId),
			Date:      optiomist.None[time.Time](),
			Hour:      optiomist.Nil[time.Time](),
			Venue:     optiomist.None[string](),
			Address:   optiomist.Nil[string](),
			Town:      optiomist.None[string](),
		}

		id, err := r.Create(event)
		if err != nil {
			t.Fatal(err)
		}
		found, err := r.Find(id)
		if err != nil {
			t.Fatal(err)
		}
		if found.Id.Value() != id {
			t.Errorf("Id = %v; expected %v", found.Id.Value(), id)
		}
		if found.AccountId.Value() != accountId {
			t.Errorf("AccountId = %v; expected %v", found.AccountId.Value(), accountId)
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
