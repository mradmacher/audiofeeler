package audiofeeler

import (
	"html/template"
	"net/http"
)

type AccountView struct {
	Id    DatabaseId
	Name  string
}

type AccountsController struct {
	app *App
	indexTemplate *template.Template
	showTemplate  *template.Template
}

func NewAccountsController(app *App) *AccountsController {
	controller := AccountsController{}
	controller.app = app

	controller.indexTemplate = template.Must(
		template.ParseFiles(
			app.templatesPath+"/accounts.gohtml",
			app.templatesPath+"/application.gohtml",
		),
	)
	controller.showTemplate = template.Must(
		template.ParseFiles(
			app.templatesPath+"/account.gohtml",
			app.templatesPath+"/application.gohtml",
		),
	)

	app.router.HandleFunc("GET /{$}", controller.homeHandler)
	app.router.HandleFunc("GET /{name}", controller.accountHandler)

	return &controller
}

func (controller *AccountsController) accountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	repo := AccountsRepo{controller.app.db}
	account, err := repo.FindByName(r.PathValue("name"))
	if err != nil {
		panic(err)
	}
	err = controller.showTemplate.ExecuteTemplate(
		w,
		"application",
		AccountView{
			Id:    account.Id,
			Name:  account.Name,
		},
	)
	if err != nil {
		panic(err)
	}
}

func (controller *AccountsController) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	repo := AccountsRepo{controller.app.db}
	accounts, err := repo.FindAll()
	if err != nil {
		panic(err)
	}
	var views []AccountView

	for _, account := range accounts {
		views = append(views, AccountView{
			Id:    account.Id,
			Name:  account.Name,
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

type App struct {
	router        *http.ServeMux
	db            *DbClient
	templatesPath string
}

func NewApp(templatesPath string, dbUrl string) (*App, error) {
	app := App{}
	app.router = http.DefaultServeMux
	app.templatesPath = templatesPath

	NewAccountsController(&app)

	var err error
	app.db, err = NewDbClient(dbUrl)
	if err != nil {
		panic(err)
	}

	return &app, nil
}

func (app *App) Start() {
	http.ListenAndServe(":3000", nil)
}

func (app *App) Cleanup() {
	app.db.Close()
}

