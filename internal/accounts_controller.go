package internal

import (
	"net/http"
)

type AccountsPresenter struct {
	Accounts []Account
}

type AccountsController struct {
	app           *App
	indexTemplate Template
	showTemplate  Template
}

func NewAccountsController(app *App) *AccountsController {
	controller := AccountsController{}
	controller.app = app

	controller.indexTemplate = app.templateEngine.MustParse("accounts")
	controller.showTemplate = app.templateEngine.MustParse("account", "account_wrapper")

	app.router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
	})

	app.router.HandleFunc("GET /accounts/{id}", func(w http.ResponseWriter, r *http.Request) {
		repo := AccountRepo{app.db}
		findResult, err := repo.Find(DatabaseId(r.PathValue("id")))
		if err != nil {
			panic(err)
		}
		if findResult.IsFound {
			err = controller.showTemplate.Execute(w, r, findResult.Record)
			if err != nil {
				panic(err)
			}
		}
	})

	app.router.HandleFunc("GET /accounts/{$}", func(w http.ResponseWriter, r *http.Request) {
		repo := AccountRepo{controller.app.db}
		accounts, err := repo.FindAll()
		if err != nil {
			panic(err)
		}
		err = controller.indexTemplate.Execute(w, r, AccountsPresenter{accounts})
		if err != nil {
			panic(err)
		}
	})

	return &controller
}
