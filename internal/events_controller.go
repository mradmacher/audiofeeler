package audiofeeler

import (
	"net/http"
)

type EventsController struct {
	app  *App
	view *EventView
}

func NewEventsController(app *App) *EventsController {
	controller := EventsController{}
	controller.app = app

	accountRepo := AccountRepo{app.db}
	eventRepo := EventRepo{app.db}

	controller.view = NewEventView(app.templateEngine)

	app.router.HandleFunc("GET /accounts/{accountName}/events/{$}", func(w http.ResponseWriter, r *http.Request) {
		accountName := r.PathValue("accountName")
		accountFindResult, err := accountRepo.FindByName(accountName)
		if err != nil {
			panic(err)
		}
		account := accountFindResult.Record
		events, err := eventRepo.FindAll(account.Id)
		if err != nil {
			panic(err)
		}
		err = controller.view.renderIndex(ViewContext{w, r}, account, events)
		if err != nil {
			panic(err)
		}
	})

	app.router.HandleFunc("GET /accounts/{accountName}/events/new", func(w http.ResponseWriter, r *http.Request) {
		accountName := r.PathValue("accountName")
		accountFindResult, err := accountRepo.FindByName(accountName)
		if err != nil {
			panic(err)
		}
		if accountFindResult.IsFound {
			err = controller.view.renderNew(ViewContext{w, r}, accountFindResult.Record, Event{})
			if err != nil {
				panic(err)
			}
		}
	})

	app.router.HandleFunc("POST /accounts/{accountName}/events", func(w http.ResponseWriter, r *http.Request) {
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
		accountFindResult, err := accountRepo.FindByName(accountName)
		if err != nil {
			panic(err)
		}
		if accountFindResult.IsFound {
			account := accountFindResult.Record
			event.AccountId = account.Id
			eventRepo.Save(event)

			http.Redirect(w, r, "/accounts/"+accountName+"/events", http.StatusFound)
		}
	})

	app.router.HandleFunc("DELETE /accounts/{accountName}/events/{eventId}", func(w http.ResponseWriter, r *http.Request) {
		accountName := r.PathValue("accountName")
		eventId := DatabaseId(r.PathValue("eventId"))
		eventRepo.Delete(eventId)

		http.Redirect(w, r, "/accounts/"+accountName+"/events", http.StatusOK)
	})

	return &controller
}
