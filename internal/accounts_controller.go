package audiofeeler

import (
	"net/http"
)

type AccountsController struct {
	app  *App
	view *AccountView
}

func NewAccountsController(app *App) *AccountsController {
	controller := AccountsController{}
	controller.app = app

	controller.view = NewAccountView(app.templateEngine)

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
		controller.view.renderShow(ViewContext{w, r}, findResult.Record)
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
	err = controller.view.renderIndex(ViewContext{w, r}, accounts)
	if err != nil {
		panic(err)
	}
}
