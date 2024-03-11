package audiofeeler

import (
	"html/template"
	"net/http"
)

type AccountView struct {
	Id uint32
	Title string
	Name string
	Url string
}

type App struct {
	router   *http.ServeMux
	indexTemplate *template.Template
	showTemplate *template.Template
}

func NewApp(templatesPath string) (*App, error) {
	app := App{}
	app.router = http.DefaultServeMux

	app.indexTemplate = template.Must(
		template.ParseFiles(
			templatesPath + "/accounts.gohtml",
			templatesPath + "/application.gohtml",
		),
	)
	app.showTemplate = template.Must(
		template.ParseFiles(
			templatesPath + "/account.gohtml",
			templatesPath + "/application.gohtml",
		),
	)
	app.MountHandlers()

	return &app, nil
}

func (app *App) MountHandlers() {
	app.router.HandleFunc("GET /{$}", app.homeHandler)
	app.router.HandleFunc("GET /{name}", app.accountHandler)
}

func (app *App) Start() {
	http.ListenAndServe(":3000", nil)
}

func (app *App) accountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := app.showTemplate.ExecuteTemplate(
		w,
		"application",
		AccountView {
			Id: uint32(1),
			Title: "Czarny Motyl",
			Name: "czarnymotyl",
			Url: "http://czarnymotyl.art",
		},
	)
	if err != nil {
		panic(err)
	}
}

func (app *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := app.indexTemplate.ExecuteTemplate(
		w,
		"application",
		struct {
			Accounts []AccountView
		}{
			[]AccountView {
				AccountView {
					Id: uint32(1),
					Title: "Czarny Motyl",
					Name: "czarnymotyl",
					Url: "http://czarnymotyl.art",
				},
				AccountView {
					Id: uint32(2),
					Title: "Karoryfer Lecolds",
					Name: "karoryfer",
					Url: "http://karoryfer.com",
				},
				AccountView {
					Id: uint32(3),
					Title: "BalkanArtz",
					Name: "balkanartz",
					Url: "http://balkanartz.eu",
				},
				AccountView {
					Id: uint32(3),
					Title: "Iglika",
					Name: "iglika",
					Url: "http://iglika.eu",
				},
			},
		},
	)
	if err != nil {
		panic(err)
	}
}
