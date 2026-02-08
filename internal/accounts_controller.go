package audiofeeler

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

	controller.indexTemplate = app.templateEngine.Parse("accounts")
	controller.showTemplate = app.templateEngine.Parse("account", "account_wrapper")

	app.router.HandleFunc("GET /{$}", controller.accountsHandler)
	app.router.HandleFunc("GET /{name}", controller.accountHandler)

	return &controller
}

func (controller *AccountsController) accountHandler(w http.ResponseWriter, r *http.Request) {
	repo := AccountRepo{controller.app.db}
	findResult, err := repo.FindByName(r.PathValue("name"))
	if err != nil {
		panic(err)
	}
	if findResult.IsFound {
		err = controller.showTemplate.Execute(w, r, findResult.Record)
		if err != nil {
			panic(err)
		}
	}
}

func (controller *AccountsController) accountsHandler(w http.ResponseWriter, r *http.Request) {
	repo := AccountRepo{controller.app.db}
	accounts, err := repo.FindAll()
	if err != nil {
		panic(err)
	}
	err = controller.indexTemplate.Execute(w, r, AccountsPresenter{accounts})
	if err != nil {
		panic(err)
	}
}
