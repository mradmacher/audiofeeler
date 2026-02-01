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

	accountRepo := AccountRepo{app.db}
	eventRepo := EventRepo{app.db}

	controller.indexTemplate = app.ParseTemplate("events", "account_wrapper")
	controller.newTemplate = app.ParseTemplate("event_form", "account_wrapper")

	app.router.HandleFunc("GET /accounts/{accountName}/events/{$}", func(w http.ResponseWriter, r *http.Request) {
		assignResponseDefaults(w)

		accountName := r.PathValue("accountName")
		account, err := accountRepo.FindByName(accountName)
		if err != nil {
			panic(err)
		}
		events, err := eventRepo.FindAll(account.Id)
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

		accountName := r.PathValue("accountName")
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

	app.router.HandleFunc("POST /accounts/{accountName}/events", func(w http.ResponseWriter, r *http.Request) {
		assignResponseDefaults(w)

		accountName := r.PathValue("accountName")
		event := Event{
			Name:        r.PostFormValue("event[name]"),
			Date:        r.PostFormValue("event[date]"),
			Hour:        r.PostFormValue("event[hour]"),
			Venue:       r.PostFormValue("event[venue]"),
			Town:        r.PostFormValue("event[town]"),
			Location:    r.PostFormValue("event[location]"),
			Description: r.PostFormValue("event[description]"),
		}
		account, err := accountRepo.FindByName(accountName)
		if err != nil {
			panic(err)
		}
		event.AccountId = account.Id
		eventRepo.Save(event)

		http.Redirect(w, r, "/accounts/"+accountName+"/events", http.StatusFound)
	})

	app.router.HandleFunc("DELETE /accounts/{accountName}/events/{eventId}", func(w http.ResponseWriter, r *http.Request) {
		assignResponseDefaults(w)

		accountName := r.PathValue("accountName")
		eventId := DatabaseId(r.PathValue("eventId"))
		eventRepo.Delete(eventId)

		http.Redirect(w, r, "/accounts/"+accountName+"/events", http.StatusOK)
	})

	return &controller
}
