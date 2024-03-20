package store

import (
	"github.com/mradmacher/audiofeeler/pkg/optiomist"
	"gotest.tools/v3/assert"
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
	assert.NilError(t, err)

	return accountId
}

func TestEventsRepo(t *testing.T) {
	teardown, db := SetupDbTest(t)
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

				assert.NilError(t, err)
				assert.Check(t, id > 0)
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
		assert.NilError(t, err)

		got, err := r.Find(id)
		assert.NilError(t, err)

		assert.Equal(t, got.Id.Value(), id)
		assert.Equal(t, got.AccountId.Value(), accountId)
		assert.Equal(t, got.Date, event.Date)
		assert.Check(t, got.Hour.IsSome() && event.Hour.IsSome() && got.Hour.Value().Format(time.TimeOnly) == event.Hour.Value().Format(time.TimeOnly))
		assert.Equal(t, got.Venue, event.Venue)
		assert.Equal(t, got.Address, event.Address)
		assert.Equal(t, got.Town, event.Town)
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
		got, err := r.Find(id)
		assert.NilError(t, err)
		assert.Equal(t, got.Id.Value(), id)
		assert.Equal(t, got.AccountId.Value(), accountId)
		assert.Check(t, got.Date.IsNone())
		assert.Check(t, got.Hour.IsNone())
		assert.Check(t, got.Venue.IsNone())
		assert.Check(t, got.Address.IsNone())
		assert.Check(t, got.Town.IsNone())
	}
}
