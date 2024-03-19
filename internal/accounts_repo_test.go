package audiofeeler

import (
	"github.com/mradmacher/audiofeeler/optiomist"
	"gotest.tools/v3/assert"
	"testing"
)

func TestAccountsRepo(t *testing.T) {
	teardown, db := setupTest(t)
	defer teardown(t)

	r := AccountsRepo{db}

	t.Run("Create with all params", testCreate_allParams(&r))
	t.Run("Create with missing params", testCreate_missingParams(&r))
	t.Run("FindAll", testFindAll(&r))
	t.Run("FindByName", testFindByName(&r))
}

func testFindAll(r *AccountsRepo) func(*testing.T) {
	return func(t *testing.T) {
		id1, err := r.Create(Account{
			Name:  optiomist.Some("account1"),
			Title: optiomist.Some("Account One"),
			Url:   optiomist.Some("http://account1.com"),
		})
		assert.NilError(t, err)

		id2, err := r.Create(Account{
			Name:  optiomist.Some("account2"),
			Title: optiomist.Some("Account Two"),
			Url:   optiomist.Some("http://account2.com"),
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

		id, err := r.Create(Account{
			Name:  optiomist.Some("someaccount"),
			Title: optiomist.Some("Some Account"),
			Url:   optiomist.Some("http://someaccount.com"),
		})
		assert.NilError(t, err)

		_, err = r.Create(Account{
			Name:  optiomist.Some("otheraccount"),
			Title: optiomist.Some("Other Account"),
			Url:   optiomist.Some("http://otheraccount.com"),
		})
		assert.NilError(t, err)

		got, err = r.FindByName("someaccount")
		assert.NilError(t, err)

		assert.Equal(t, got.Id.Value(), id)
		assert.Equal(t, got.Name.Value(), "someaccount")
		assert.Equal(t, got.Title.Value(), "Some Account")
		assert.Equal(t, got.Url.Value(), "http://someaccount.com")

		got, err = r.FindByName("yetanotheraccount")
		assert.ErrorIs(t, err, newRecordNotFound())
	}
}

func testCreate_missingParams(r *AccountsRepo) func(*testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name    string
			account Account
		}{
			{
				"missing url",
				Account{
					Name:  optiomist.Some("onlyexample"),
					Title: optiomist.Some("Example"),
					Url:   optiomist.None[string](),
				},
			}, {
				"missing name",
				Account{
					Name:  optiomist.None[string](),
					Title: optiomist.Some("Example"),
					Url:   optiomist.Some("http://onlyexample.com"),
				},
			}, {
				"missing title",
				Account{
					Name: optiomist.None[string](),
					Url:  optiomist.Some("http://onlyexample.com"),
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				_, err := r.Create(test.account)

				if err == nil {
					t.Error("It should not create record with missing data")
				}
			})
		}
	}
}

func testCreate_allParams(r *AccountsRepo) func(*testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name    string
			account Account
		}{
			{
				"all params",
				Account{
					Name:  optiomist.Some("example"),
					Title: optiomist.Some("Example"),
					Url:   optiomist.Some("http://example.com"),
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
