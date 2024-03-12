package audiofeeler

import (
	"html/template"
	"net/http"
	"os"
)

type AccountView struct {
	Id    uint32
	Title string
	Name  string
	Url   string
}

type App struct {
	router        *http.ServeMux
	indexTemplate *template.Template
	showTemplate  *template.Template
	db            *DbClient
}

func NewApp(templatesPath string) (*App, error) {
	app := App{}
	app.router = http.DefaultServeMux

	app.indexTemplate = template.Must(
		template.ParseFiles(
			templatesPath+"/accounts.gohtml",
			templatesPath+"/application.gohtml",
		),
	)
	app.showTemplate = template.Must(
		template.ParseFiles(
			templatesPath+"/account.gohtml",
			templatesPath+"/application.gohtml",
		),
	)
	app.MountHandlers()

	var err error
	app.db, err = NewDbClient(os.Getenv("AUDIOFEELER_DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	return &app, nil
}

func (app *App) MountHandlers() {
	app.router.HandleFunc("GET /{$}", app.homeHandler)
	app.router.HandleFunc("GET /{name}", app.accountHandler)
}

func (app *App) Start() {
	http.ListenAndServe(":3000", nil)
}

func (app *App) Cleanup() {
	app.db.Close()
}

func (app *App) accountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	repo := AccountsRepo{app.db}
	account, err := repo.FindByName(r.PathValue("name"))
	if err != nil {
		panic(err)
	}
	err = app.showTemplate.ExecuteTemplate(
		w,
		"application",
		AccountView{
			Id:    account.Id.Value(),
			Title: account.Title.Value(),
			Name:  account.Name.Value(),
			Url:   account.Url.Value(),
		},
	)
	if err != nil {
		panic(err)
	}
}

func (app *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	repo := AccountsRepo{app.db}
	accounts, err := repo.FindAll()
	if err != nil {
		panic(err)
	}
	var views []AccountView

	for _, account := range accounts {
		views = append(views, AccountView{
			Id:    account.Id.Value(),
			Title: account.Title.Value(),
			Name:  account.Name.Value(),
			Url:   account.Url.Value(),
		})
	}
	err = app.indexTemplate.ExecuteTemplate(
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
