package internal

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestAccountRepo(t *testing.T) {
	teardown, db := setupDbTest(t)
	defer teardown(t)

	r := AccountRepo{db}

	t.Run("FindAll", testFindAll(&r))
	t.Run("Find", testFind(&r))
	t.Run("FindByName", testFindByName(&r))
	t.Run("Create", testCreate(&r))
}

func testFindAll(r *AccountRepo) func(*testing.T) {
	return func(t *testing.T) {
		id1, err := r.Save(Account{
			Name: "account1",
		})
		assert.NilError(t, err)

		id2, err := r.Save(Account{
			Name: "account2",
		})
		assert.NilError(t, err)

		accounts, err := r.FindAll()
		assert.NilError(t, err)
		t.Log(accounts)
		t.Log(id1, id2)
		assert.Equal(t, 2, len(accounts))
	}
}

func testCreate(r *AccountRepo) func(*testing.T) {
	return func(t *testing.T) {
		account := Account{
			Name: "example",
		}

		id, err := r.Save(account)

		assert.NilError(t, err)
		assert.Check(t, IsDatabaseIdSet(id))
	}
}

func testFind(r *AccountRepo) func(*testing.T) {
	return func(t *testing.T) {
		var got Account

		result, err := r.Find("someaccount")
		assert.NilError(t, err)
		assert.Check(t, !result.IsFound)

		id, err := r.Save(Account{
			Name: "Test Account",
		})
		assert.NilError(t, err)

		result, err = r.Find(id)
		assert.NilError(t, err)
		assert.Check(t, result.IsFound)
		got = result.Record

		assert.Equal(t, got.Id, id)
		assert.Equal(t, got.Name, "Test Account")
	}
}

func testFindByName(r *AccountRepo) func(*testing.T) {
	return func(t *testing.T) {
		var got Account

		gotResult, err := r.FindByName("someaccount")
		assert.NilError(t, err)
		assert.Check(t, !gotResult.IsFound)

		id, err := r.Save(Account{
			Name: "someaccount",
		})
		assert.NilError(t, err)

		_, err = r.Save(Account{
			Name: "otheraccount",
		})
		assert.NilError(t, err)

		gotResult, err = r.FindByName("someaccount")
		assert.NilError(t, err)
		assert.Check(t, gotResult.IsFound)
		got = gotResult.Record

		assert.Equal(t, got.Id, id)
		assert.Equal(t, got.Name, "someaccount")

		gotResult, err = r.FindByName("yetanotheraccount")
		assert.NilError(t, err)
		assert.Check(t, !gotResult.IsFound)
	}
}
