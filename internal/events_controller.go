package audiofeeler

import (
	"html/template"
	"net/http"
)

type EventsController struct {
	app           *App
	indexTemplate *template.Template
	newTemplate   *template.Template
}

func NewEventsController(app *App) *EventsController {
	controller := EventsController{}
	controller.app = app

	controller.indexTemplate = app.ParseTemplate("events", "account_wrapper")
	controller.newTemplate = app.ParseTemplate("event_form", "account_wrapper")

	app.router.HandleFunc("GET /accounts/{accountName}/events/{$}", func(w http.ResponseWriter, r *http.Request) {
		assignResponseDefaults(w)

		accountRepo := AccountRepo{controller.app.db}
		repo := EventRepo{controller.app.db}
		accountName := r.PathValue("accountName")
		//id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		account, err := accountRepo.FindByName(accountName)
		if err != nil {
			panic(err)
		}
		events, err := repo.FindAll(account.Id)
		if err != nil {
			panic(err)
		}
		err = controller.indexTemplate.ExecuteTemplate(
			w,
			selectTemplateName(r),
			struct {
				Account Account
				Events  []Event
			}{
				account,
				events,
			},
		)
		if err != nil {
			panic(err)
		}
	})

	app.router.HandleFunc("GET /accounts/{accountName}/events/new", func(w http.ResponseWriter, r *http.Request) {
		assignResponseDefaults(w)

		accountRepo := AccountRepo{controller.app.db}
		accountName := r.PathValue("accountName")
		//id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		account, err := accountRepo.FindByName(accountName)
		if err != nil {
			panic(err)
		}
		err = controller.newTemplate.ExecuteTemplate(
			w,
			selectTemplateName(r),
			struct {
				Account Account
				Event   Event
			}{
				account,
				Event{},
			},
		)
		if err != nil {
			panic(err)
		}
	})

	return &controller
}
