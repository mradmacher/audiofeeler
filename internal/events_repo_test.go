package audiofeeler

import (
	. "github.com/mradmacher/audiofeeler/pkg/optiomist"
	"gotest.tools/v3/assert"
	"testing"
)

func setupAccount(db *DbClient, t *testing.T) int64 {
	accountsRepo := AccountsRepo{db}
	accountId, err := accountsRepo.Create(AccountParams{
		Name:  Some("example"),
		SourceDir: Some("/here"),
	})
	assert.NilError(t, err)

	return accountId
}

func TestEventsRepo(t *testing.T) {
	teardown, db := setupDbTest(t)
	defer teardown(t)

	accountId := setupAccount(db, t)

	r := EventsRepo{db}

	t.Run("Create", testEventsRepo_Create(&r, accountId))
	t.Run("Find not nil values", testEventsRepo_Find_not_nils(&r, accountId))
	t.Run("Find nil values", testEventsRepo_Find_nils(&r, accountId))
}

func testEventsRepo_Create(r *EventsRepo, accountId int64) func(*testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name  string
			event EventParams
		}{
			{
				"some params",
				EventParams{
					AccountId: Some(accountId),
					Date:      Some("2024-02-01"),
					Hour:      Some("21:00"),
					Venue:     Some("Some venue"),
					Address:   Some("Some address"),
					City:      Some("Some city"),
					Place:     Some("Some place"),
				},
			}, {
				"none params",
				EventParams{
					AccountId: Some(accountId),
					Date:      None[string](),
					Hour:      None[string](),
					Venue:     None[string](),
					Address:   None[string](),
					City:      None[string](),
					Place:     None[string](),
				},
			}, {
				"nil params",
				EventParams{
					AccountId: Some(accountId),
					Date:      Some(""),
					Hour:      Some(""),
					Venue:     Some(""),
					Address:   Some(""),
					City:      Some(""),
					Place:     Some(""),
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

func testEventsRepo_Find_not_nils(r *EventsRepo, accountId int64) func(*testing.T) {
	return func(t *testing.T) {
		event := EventParams{
			AccountId: Some(accountId),
			Date:      Some("2024-02-02"),
			Hour:      Some("21:00"),
			Venue:     Some("Some venue"),
			Address:   Some("Some address"),
			City:      Some("Some city"),
			Place:     Some("Some place"),
		}

		id, err := r.Create(event)
		assert.NilError(t, err)

		got, err := r.Find(id)
		assert.NilError(t, err)

		assert.Equal(t, got.Id, id)
		assert.Equal(t, got.AccountId, accountId)
		assert.Equal(t, got.Date, event.Date)
		assert.Check(t, got.Hour, event.Hour)
		assert.Equal(t, got.Venue, event.Venue)
		assert.Equal(t, got.Address, event.Address)
		assert.Equal(t, got.City, event.City)
		assert.Equal(t, got.Place, event.Place)
	}
}

func testEventsRepo_Find_nils(r *EventsRepo, accountId int64) func(*testing.T) {
	return func(t *testing.T) {
		event := EventParams{
			AccountId: Some(accountId),
			Date:      None[string](),
			Hour:      Some(""),
			Venue:     None[string](),
			Address:   Some(""),
			City:      None[string](),
			Place:     Some(""),
		}

		id, err := r.Create(event)
		if err != nil {
			t.Fatal(err)
		}
		got, err := r.Find(id)
		assert.NilError(t, err)
		assert.Equal(t, got.Id, id)
		assert.Equal(t, got.AccountId, accountId)
		assert.Equal(t, got.Date, "")
		assert.Equal(t, got.Hour, "")
		assert.Equal(t, got.Venue, "")
		assert.Equal(t, got.Address, "")
		assert.Equal(t, got.City, "")
		assert.Equal(t, got.Place, "")
	}
}
