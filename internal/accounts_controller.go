package audiofeeler

import (
	"html/template"
	"net/http"
)

type AccountView struct {
	Id   DatabaseId
	Name string
}

type AccountsController struct {
	app           *App
	indexTemplate *template.Template
	showTemplate  *template.Template
}

func NewAccountsController(app *App) *AccountsController {
	controller := AccountsController{}
	controller.app = app

	controller.indexTemplate = app.ParseTemplate("accounts")
	controller.showTemplate = app.ParseTemplate("account", "account_wrapper")

	app.router.HandleFunc("GET /{$}", controller.accountsHandler)
	app.router.HandleFunc("GET /{name}", controller.accountHandler)

	return &controller
}

func (controller *AccountsController) accountHandler(w http.ResponseWriter, r *http.Request) {
	assignResponseDefaults(w)
	repo := AccountsRepo{controller.app.db}
	account, err := repo.FindByName(r.PathValue("name"))
	if err != nil {
		panic(err)
	}
	err = controller.showTemplate.ExecuteTemplate(
		w,
		"application",
		AccountView{
			Id:   account.Id,
			Name: account.Name,
		},
	)
	if err != nil {
		panic(err)
	}
}

func (controller *AccountsController) accountsHandler(w http.ResponseWriter, r *http.Request) {
	assignResponseDefaults(w)
	repo := AccountsRepo{controller.app.db}
	accounts, err := repo.FindAll()
	if err != nil {
		panic(err)
	}
	var views []AccountView

	for _, account := range accounts {
		views = append(views, AccountView{
			Id:   account.Id,
			Name: account.Name,
		})
	}
	err = controller.indexTemplate.ExecuteTemplate(
		w,
		"application",
		struct {
			Accounts []AccountView
		}{
			views,
		},
	)
	if err != nil {
		panic(err)
	}
}
