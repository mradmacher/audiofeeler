package audiofeeler

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestAccountRepo(t *testing.T) {
	teardown, db := setupDbTest(t)
	defer teardown(t)

	r := AccountRepo{db}

	t.Run("Create with all params", testCreate_allParams(&r))
	t.Run("Create with missing params", testCreate_missingParams(&r))
	t.Run("Create with duplicated name", testCreate_duplicatedName(&r))
	t.Run("FindAll", testFindAll(&r))
	t.Run("FindByName", testFindByName(&r))
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

func testCreate_duplicatedName(r *AccountRepo) func(*testing.T) {
	return func(t *testing.T) {
		account := Account{
			Name: "this-is-unique",
		}
		_, err := r.Save(account)
		assert.NilError(t, err)

		dupAccount := Account{
			Name: "this-is-unique",
		}
		_, err = r.Save(dupAccount)
		assert.Check(t, err != nil, "It should not create record with duplicated name")
	}
}

func testCreate_missingParams(r *AccountRepo) func(*testing.T) {
	return func(t *testing.T) {
		account := Account{
			Name: "",
		}

		_, err := r.Save(account)
		assert.Check(t, err != nil, "It should not create record with missing data")
	}
}

func testCreate_allParams(r *AccountRepo) func(*testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name    string
			account Account
		}{
			{
				"all params",
				Account{
					Name: "example",
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				id, err := r.Save(test.account)

				assert.NilError(t, err)
				assert.Check(t, IsDatabaseIdSet(id))
			})
		}
	}
}

func testFindByName(r *AccountRepo) func(*testing.T) {
	return func(t *testing.T) {
		var got Account

		got, err := r.FindByName("someaccount")
		assert.ErrorIs(t, err, newRecordNotFound())

		id, err := r.Save(Account{
			Name: "someaccount",
		})
		assert.NilError(t, err)

		_, err = r.Save(Account{
			Name: "otheraccount",
		})
		assert.NilError(t, err)

		got, err = r.FindByName("someaccount")
		assert.NilError(t, err)

		assert.Equal(t, got.Id, id)
		assert.Equal(t, got.Name, "someaccount")

		got, err = r.FindByName("yetanotheraccount")
		assert.ErrorIs(t, err, newRecordNotFound())
	}
}
