package audiofeeler

import (
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
)

type App struct {
	router   *chi.Mux
	template *template.Template
}

func NewApp(templatesPath string) (*App, error) {
	app := App{}
	app.router = chi.NewRouter()

	var err error

	app.template, err = template.ParseFiles(templatesPath + "/index.html")
	if err != nil {
		return nil, err
	}
	app.MountHandlers()

	return &app, nil
}

func (app *App) MountHandlers() {
	app.router.Get("/", app.homeHandler)
}

func (app *App) Start() {
	http.ListenAndServe(":3000", app.router)
}

func (app *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := app.template.Execute(w, struct{}{})
	if err != nil {
		panic(err)
	}
}
