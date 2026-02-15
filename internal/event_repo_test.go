package internal

import (
	"gotest.tools/v3/assert"
	"testing"
)

func setupAccount(db *DbClient, t *testing.T) DatabaseId {
	accountRepo := AccountRepo{db}
	accountId, err := accountRepo.Save(
		Account{
			Name: "example",
		},
	)
	assert.NilError(t, err)

	return accountId
}

func TestEventRepo(t *testing.T) {
	teardown, db := setupDbTest(t)
	defer teardown(t)

	accountId := setupAccount(db, t)

	r := EventRepo{db}

	t.Run("Save create", testEventRepo_Save_create(&r, accountId))
	t.Run("Save update", testEventRepo_Save_update(&r, accountId))
	t.Run("Delete", testEventRepo_Delete(&r, accountId))
	t.Run("Find not nil values", testEventRepo_Find_not_nils(&r, DatabaseId(accountId)))
	t.Run("Find nil values", testEventRepo_Find_nils(&r, DatabaseId(accountId)))
}

func testEventRepo_Save_create(r *EventRepo, accountId DatabaseId) func(*testing.T) {
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
					Status:      CurrentEvent,
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
					Status:      CurrentEvent,
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				id, err := r.Save(test.event)

				assert.NilError(t, err)
				assert.Check(t, IsDatabaseIdSet(id))
			})
		}
	}
}

func testEventRepo_Save_update(r *EventRepo, accountId DatabaseId) func(*testing.T) {
	return func(t *testing.T) {
		event := Event{
			AccountId:   DatabaseId(accountId),
			Date:        "2024-02-01",
			Hour:        "21:00",
			Name:        "Some festival",
			Venue:       "Some venue",
			Location:    "Some location",
			Town:        "Some town",
			Description: "Some description",
			Status:      CurrentEvent,
		}
		newEvent := Event{
			AccountId:   DatabaseId(accountId),
			Date:        "2025-03-02",
			Hour:        "22:00",
			Name:        "Some new festival",
			Venue:       "Some new venue",
			Location:    "Some new location",
			Town:        "Some new town",
			Description: "Some new description",
			Status:      ArchivedEvent,
		}

		t.Run("updates existing event", func(t *testing.T) {
			var got Event
			id, err := r.Save(event)
			assert.NilError(t, err)

			newEvent.Id = id
			id, err = r.Save(newEvent)
			assert.NilError(t, err)

			gotResult, err := r.Find(id)
			assert.NilError(t, err)
			assert.Check(t, gotResult.IsFound)
			got = gotResult.Record

			assert.Equal(t, got.Id, newEvent.Id)
			assert.Equal(t, got.AccountId, newEvent.AccountId)
			assert.Equal(t, got.Date, newEvent.Date)
			assert.Equal(t, got.Name, newEvent.Name)
			assert.Equal(t, got.Hour, newEvent.Hour)
			assert.Equal(t, got.Venue, newEvent.Venue)
			assert.Equal(t, got.Town, newEvent.Town)
			assert.Equal(t, got.Location, newEvent.Location)
			assert.Equal(t, got.Description, newEvent.Description)
			assert.Equal(t, got.Status, newEvent.Status)
		})
	}
}

func testEventRepo_Delete(r *EventRepo, accountId DatabaseId) func(*testing.T) {
	return func(t *testing.T) {
		event := Event{
			AccountId:   DatabaseId(accountId),
			Date:        "2024-02-01",
			Hour:        "21:00",
			Name:        "Some festival",
			Venue:       "Some venue",
			Location:    "Some location",
			Town:        "Some town",
			Description: "Some description",
			Status:      CurrentEvent,
		}

		t.Run("deletes existing event", func(t *testing.T) {
			id, err := r.Save(event)
			assert.NilError(t, err)
			var gotResult FindResult[Event]
			gotResult, err = r.Find(id)
			assert.NilError(t, err)
			assert.Check(t, gotResult.IsFound)

			r.Delete(id)
			gotResult, err = r.Find(id)
			assert.NilError(t, err)
			assert.Check(t, !gotResult.IsFound)
		})
	}
}
func testEventRepo_Find_not_nils(r *EventRepo, accountId DatabaseId) func(*testing.T) {
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
			Status:      ArchivedEvent,
		}

		var got Event
		id, err := r.Save(event)
		assert.NilError(t, err)

		gotResult, err := r.Find(id)
		assert.NilError(t, err)
		assert.Check(t, gotResult.IsFound)
		got = gotResult.Record

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

func testEventRepo_Find_nils(r *EventRepo, accountId DatabaseId) func(*testing.T) {
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
			Status:      ArchivedEvent,
		}

		id, err := r.Save(event)
		if err != nil {
			t.Fatal(err)
		}
		gotResult, err := r.Find(id)
		assert.NilError(t, err)
		assert.Check(t, gotResult.IsFound)
		got := gotResult.Record

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
