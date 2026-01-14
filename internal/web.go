package audiofeeler

import (
	"html/template"
	"net/http"
)

func assignResponseDefaults(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func selectTemplateName(r *http.Request) string {
	templateName := "application"
	if r.Header.Get("Hx-Request") == "true" {
		templateName = "content"
	}
	return templateName
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

	var err error
	app.db, err = NewDbClient(dbUrl)
	if err != nil {
		panic(err)
	}

	NewAccountsController(&app)
	NewEventsController(&app)

	return &app, nil
}

func (app *App) ParseTemplate(names ...string) *template.Template {
	files := []string{app.templatesPath + "/application.gohtml"}
	for _, name := range names {
		files = append(files, app.templatesPath+"/"+name+".gohtml")
	}
	return template.Must(template.ParseFiles(files...))
}

func (app *App) Start() {
	http.ListenAndServe(":3000", nil)
}

func (app *App) Cleanup() {
	app.db.Close()
}
