package audiofeeler

import (
	"gotest.tools/v3/assert"
	"testing"
)

func setupAccount(db *DbClient, t *testing.T) DatabaseId {
	accountsRepo := AccountsRepo{db}
	accountId, err := accountsRepo.Save(
		Account{
			Name: "example",
		},
	)
	assert.NilError(t, err)

	return accountId
}

func TestEventsRepo(t *testing.T) {
	teardown, db := setupDbTest(t)
	defer teardown(t)

	accountId := setupAccount(db, t)

	r := EventsRepo{db}

	t.Run("Save", testEventsRepo_Save(&r, accountId))
	t.Run("Find not nil values", testEventsRepo_Find_not_nils(&r, DatabaseId(accountId)))
	t.Run("Find nil values", testEventsRepo_Find_nils(&r, DatabaseId(accountId)))
}

func testEventsRepo_Save(r *EventsRepo, accountId DatabaseId) func(*testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name  string
			event Event
		}{
			{
				"some params",
				Event{
					AccountId:   DatabaseId(accountId),
					Date:        "2024-02-01",
					Hour:        "21:00",
					Name:        "Some festival",
					Venue:       "Some venue",
					Location:    "Some location",
					Town:        "Some town",
					Description: "Some description",
					Status:      EventCurrent,
				},
			}, {
				"none params",
				Event{
					AccountId:   DatabaseId(accountId),
					Name:        "",
					Date:        "",
					Hour:        "",
					Venue:       "",
					Town:        "",
					Location:    "",
					Description: "",
					Status:      EventCurrent,
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				id, err := r.Save(test.event)

				assert.NilError(t, err)
				assert.Check(t, id > 0)
			})
		}
	}
}

func testEventsRepo_Find_not_nils(r *EventsRepo, accountId DatabaseId) func(*testing.T) {
	return func(t *testing.T) {
		event := Event{
			AccountId:   accountId,
			Name:        "Event",
			Date:        "2024-02-02",
			Hour:        "21:00",
			Venue:       "Some venue",
			Town:        "Some town",
			Location:    "Some location",
			Description: "Some description",
			Status:      EventArchived,
		}

		id, err := r.Save(event)
		assert.NilError(t, err)

		got, err := r.Find(id)
		assert.NilError(t, err)

		assert.Equal(t, got.Id, id)
		assert.Equal(t, got.AccountId, accountId)
		assert.Equal(t, got.Date, event.Date)
		assert.Equal(t, got.Name, event.Name)
		assert.Equal(t, got.Hour, event.Hour)
		assert.Equal(t, got.Venue, event.Venue)
		assert.Equal(t, got.Town, event.Town)
		assert.Equal(t, got.Location, event.Location)
		assert.Equal(t, got.Description, event.Description)
		assert.Equal(t, got.Status, event.Status)
	}
}

func testEventsRepo_Find_nils(r *EventsRepo, accountId DatabaseId) func(*testing.T) {
	return func(t *testing.T) {
		event := Event{
			AccountId:   accountId,
			Name:        "",
			Date:        "",
			Hour:        "",
			Venue:       "",
			Town:        "",
			Location:    "",
			Description: "",
			Status:      EventArchived,
		}

		id, err := r.Save(event)
		if err != nil {
			t.Fatal(err)
		}
		got, err := r.Find(id)
		assert.NilError(t, err)
		assert.Equal(t, got.Id, id)
		assert.Equal(t, got.AccountId, accountId)
		assert.Equal(t, got.Name, "")
		assert.Equal(t, got.Date, "")
		assert.Equal(t, got.Hour, "")
		assert.Equal(t, got.Venue, "")
		assert.Equal(t, got.Town, "")
		assert.Equal(t, got.Location, "")
		assert.Equal(t, got.Description, "")
	}
}
