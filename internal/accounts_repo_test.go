package audiofeeler

import (
	. "github.com/mradmacher/audiofeeler/pkg/optiomist"
	"gotest.tools/v3/assert"
	"testing"
)

func TestAccountsRepo(t *testing.T) {
	teardown, db := setupDbTest(t)
	defer teardown(t)

	r := AccountsRepo{db}

	t.Run("Create with all params", testCreate_allParams(&r))
	t.Run("Create with missing params", testCreate_missingParams(&r))
	t.Run("Create with duplicated name", testCreate_duplicatedName(&r))
	t.Run("FindAll", testFindAll(&r))
	t.Run("FindByName", testFindByName(&r))
}

func testFindAll(r *AccountsRepo) func(*testing.T) {
	return func(t *testing.T) {
		id1, err := r.Create(AccountParams{
			Name:  Some("account1"),
			SourceDir:   Some("/here"),
		})
		assert.NilError(t, err)

		id2, err := r.Create(AccountParams{
			Name:  Some("account2"),
			SourceDir: Some("/there"),
		})
		assert.NilError(t, err)

		accounts, err := r.FindAll()
		assert.NilError(t, err)
		t.Log(accounts)
		t.Log(id1, id2)
	}
}

func testFindByName(r *AccountsRepo) func(*testing.T) {
	return func(t *testing.T) {
		var got Account

		got, err := r.FindByName("someaccount")
		assert.ErrorIs(t, err, newRecordNotFound())

		id, err := r.Create(AccountParams{
			Name:  Some("someaccount"),
			SourceDir:  Some("/here"),
		})
		assert.NilError(t, err)

		_, err = r.Create(AccountParams{
			Name:  Some("otheraccount"),
			SourceDir:   Some("/there"),
		})
		assert.NilError(t, err)

		got, err = r.FindByName("someaccount")
		assert.NilError(t, err)

		assert.Equal(t, got.Id, id)
		assert.Equal(t, got.Name, "someaccount")
		assert.Equal(t, got.SourceDir, "/here")

		got, err = r.FindByName("yetanotheraccount")
		assert.ErrorIs(t, err, newRecordNotFound())
	}
}

func testCreate_duplicatedName(r *AccountsRepo) func(*testing.T) {
	return func(t *testing.T) {
		account := AccountParams{
			Name:  Some("this-is-unique"),
			SourceDir: Some("/here"),
		}
		_, err := r.Create(account)
		assert.NilError(t, err)

		dupAccount := AccountParams{
			Name:  Some("this-is-unique"),
			SourceDir: Some("/there"),
		}
		_, err = r.Create(dupAccount)
		assert.Check(t, err != nil, "It should not create record with duplicated name")
	}
}

func testCreate_missingParams(r *AccountsRepo) func(*testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name    string
			account AccountParams
		}{
			{
				"empty name",
				AccountParams{
					Name:  Some(""),
					SourceDir:  Some("/here"),
				},
			}, {
				"missing name",
				AccountParams{
					Name: None[string](),
					SourceDir:  Some("/here"),
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				_, err := r.Create(test.account)

				assert.Check(t, err != nil, "It should not create record with missing data")
			})
		}
	}
}

func testCreate_allParams(r *AccountsRepo) func(*testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name    string
			account AccountParams
		}{
			{
				"all params",
				AccountParams{
					Name:  Some("example"),
					SourceDir:   Some("/here"),
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				id, err := r.Create(test.account)

				assert.NilError(t, err)
				assert.Check(t, id > 0)
			})
		}
	}
}
