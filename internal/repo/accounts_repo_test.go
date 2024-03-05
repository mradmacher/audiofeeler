package repo

import (
	"github.com/mradmacher/audiofeeler/internal"
	"github.com/mradmacher/audiofeeler/optiomist"
	"testing"
)

func TestAccountsRepo(t *testing.T) {
	teardown, db := setupTest(t)
	defer teardown(t)

	r := AccountsRepo{db}

	t.Run("Create with all params", testCreate_allParams(&r))
	t.Run("Create with missing params", testCreate_missingParams(&r))
	t.Run("FindAll", testFindAll(&r))
}

func testFindAll(r *AccountsRepo) func(*testing.T) {
	return func(t *testing.T) {
		id1, err := r.Create(audiofeeler.Account{
			Name: optiomist.Some("account1"),
			Url:  optiomist.Some("http://account1.com"),
		})
		if err != nil {
			t.Fatalf("Failed to create account: %v", err)
		}

		id2, err := r.Create(audiofeeler.Account{
			Name: optiomist.Some("account2"),
			Url:  optiomist.Some("http://account2.com"),
		})
		if err != nil {
			t.Fatalf("Failed to create account: %v", err)
		}

		accounts, err := r.FindAll()
		if err != nil {
			t.Fatalf("Error retrieving accounts: %v", err)
		}
		t.Log(accounts)
		t.Log(id1, id2)
	}
}

func testCreate_missingParams(r *AccountsRepo) func(*testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name    string
			account audiofeeler.Account
		}{
			{
				"missing url",
				audiofeeler.Account{
					Name: optiomist.Some("onlyexample"),
					Url:  optiomist.None[string](),
				},
			}, {
				"missing name",
				audiofeeler.Account{
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
			account audiofeeler.Account
		}{
			{
				"all params",
				audiofeeler.Account{
					Name: optiomist.Some("example"),
					Url:  optiomist.Some("http://example.com"),
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Log(test.account)
				id, err := r.Create(test.account)

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
