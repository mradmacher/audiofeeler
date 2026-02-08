package audiofeeler

import (
	"net/http"
)

type App struct {
	router         *http.ServeMux
	db             *DbClient
	templateEngine TemplateEngine
}

func NewApp(templateEngine TemplateEngine, dbUrl string) (*App, error) {
	app := App{}
	app.router = http.DefaultServeMux
	app.templateEngine = templateEngine

	var err error
	app.db, err = NewDbClient(dbUrl)
	if err != nil {
		panic(err)
	}

	NewAccountsController(&app)
	NewEventsController(&app)

	return &app, nil
}

func (app *App) Start() {
	http.ListenAndServe(":3000", nil)
}

func (app *App) Cleanup() {
	app.db.Close()
}
