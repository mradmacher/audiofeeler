package audiofeeler

import (
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
	t.Run("FindByName", testFindByName(&r))
}

func testFindAll(r *AccountsRepo) func(*testing.T) {
	return func(t *testing.T) {
		id1, err := r.Create(Account{
			Name:  optiomist.Some("account1"),
			Title: optiomist.Some("Account One"),
			Url:   optiomist.Some("http://account1.com"),
		})
		if err != nil {
			t.Fatalf("Failed to create account: %v", err)
		}

		id2, err := r.Create(Account{
			Name:  optiomist.Some("account2"),
			Title: optiomist.Some("Account Two"),
			Url:   optiomist.Some("http://account2.com"),
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

func testFindByName(r *AccountsRepo) func(*testing.T) {
	return func(t *testing.T) {
		var got Account

		got, err := r.FindByName("someaccount")
		if _, ok := err.(RecordNotFound); !ok {
			t.Errorf("error == %v; expected: RecordNotFound", err)
		}

		id, err := r.Create(Account{
			Name:  optiomist.Some("someaccount"),
			Title: optiomist.Some("Some Account"),
			Url:   optiomist.Some("http://someaccount.com"),
		})
		if err != nil {
			t.Fatalf("Failed to create account: %v", err)
		}

		_, err = r.Create(Account{
			Name:  optiomist.Some("otheraccount"),
			Title: optiomist.Some("Other Account"),
			Url:   optiomist.Some("http://otheraccount.com"),
		})
		if err != nil {
			t.Fatalf("Failed to create account: %v", err)
		}

		got, err = r.FindByName("someaccount")
		if err != nil {
			if _, ok := err.(RecordNotFound); ok {
				t.Error("got RecordNotFound; expected some data")
			} else {
				t.Fatalf("Error retrieving account: %v", err)
			}
		}
		if got.Id.Value() != id {
			t.Errorf("got == %v; expected: %v", got.Id.Value(), id)
		}
		if got.Name.Value() != "someaccount" {
			t.Errorf("got == %v; expected: %v", got.Name.Value(), "someaccount")
		}
		if got.Title.Value() != "Some Account" {
			t.Errorf("got == %v; expected: %v", got.Title.Value(), "Some Account")
		}
		if got.Url.Value() != "http://someaccount.com" {
			t.Errorf("got == %v; expected: %v", got.Url.Value(), "http://someaccount.com")
		}

		got, err = r.FindByName("yetanotheraccount")
		if _, ok := err.(RecordNotFound); !ok {
			t.Errorf("error == %v; expected: RecordNotFound", err)
		}
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
