package internal

import (
	"net/http"
)

type App struct {
	router         *http.ServeMux
	db             DbEngine
	templateEngine TemplateEngine
}

func NewApp(templateEngine TemplateEngine, dbEngine DbEngine) *App {
	app := App{}
	app.router = http.DefaultServeMux
	app.templateEngine = templateEngine
	app.db = dbEngine

	NewAccountsController(&app)
	NewEventsController(&app)

	return &app
}

func (app *App) Start() {
	http.ListenAndServe(":3000", nil)
}

func (app *App) Cleanup() {
	app.db.Close()
}
