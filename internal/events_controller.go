package internal

import (
	"net/http"
)

type EventPresenter struct {
	Account Account
	Event   Event
}

type EventsPresenter struct {
	Account Account
	Events  []Event
}

type EventsController struct {
	app           *App
	indexTemplate Template
	formTemplate  Template
}

func NewEventsController(app *App) *EventsController {
	controller := EventsController{}
	controller.app = app

	accountRepo := AccountRepo{app.db}
	eventRepo := EventRepo{app.db}

	controller.indexTemplate = app.templateEngine.MustParse("events", "account_wrapper")
	controller.formTemplate = app.templateEngine.MustParse("event_form", "account_wrapper")

	app.router.HandleFunc("GET /accounts/{accountId}/events/{$}", func(w http.ResponseWriter, r *http.Request) {
		accountFindResult, err := accountRepo.Find(DatabaseId(r.PathValue("accountId")))
		if err != nil {
			panic(err)
		}
		account := accountFindResult.Record
		events, err := eventRepo.FindAll(account.Id)
		if err != nil {
			panic(err)
		}
		err = controller.indexTemplate.Execute(w, r, EventsPresenter{account, events})
		if err != nil {
			panic(err)
		}
	})

	app.router.HandleFunc("GET /accounts/{accountId}/events/new", func(w http.ResponseWriter, r *http.Request) {
		accountFindResult, err := accountRepo.Find(DatabaseId(r.PathValue("accountId")))
		if err != nil {
			panic(err)
		}
		if accountFindResult.IsFound {
			err = controller.formTemplate.Execute(w, r, EventPresenter{accountFindResult.Record, Event{}})
			if err != nil {
				panic(err)
			}
		}
	})

	app.router.HandleFunc("GET /accounts/{accountId}/events/{eventId}/edit", func(w http.ResponseWriter, r *http.Request) {
		accountFindResult, err := accountRepo.Find(DatabaseId(r.PathValue("accountId")))
		if err != nil {
			panic(err)
		}
		if accountFindResult.IsFound {
			account := accountFindResult.Record
			eventId := DatabaseId(r.PathValue("eventId"))
			eventFindResult, err := eventRepo.Find(eventId)
			if err != nil {
				panic(err)
			}
			if eventFindResult.IsFound {
				event := eventFindResult.Record
				err = controller.formTemplate.Execute(w, r, EventPresenter{account, event})
				if err != nil {
					panic(err)
				}
			}
		}
	})

	app.router.HandleFunc("POST /accounts/{accountId}/events", func(w http.ResponseWriter, r *http.Request) {
		accountId := r.PathValue("accountId")
		event := Event{
			Name:        r.PostFormValue("event[name]"),
			Date:        r.PostFormValue("event[date]"),
			Hour:        r.PostFormValue("event[hour]"),
			Venue:       r.PostFormValue("event[venue]"),
			Town:        r.PostFormValue("event[town]"),
			Location:    r.PostFormValue("event[location]"),
			Description: r.PostFormValue("event[description]"),
		}
		accountFindResult, err := accountRepo.Find(DatabaseId(accountId))
		if err != nil {
			panic(err)
		}
		if accountFindResult.IsFound {
			account := accountFindResult.Record
			event.AccountId = account.Id
			eventRepo.Save(event)

			http.Redirect(w, r, "/accounts/"+accountId+"/events", http.StatusSeeOther)
		}
	})

	app.router.HandleFunc("PUT /accounts/{accountId}/events/{eventId}", func(w http.ResponseWriter, r *http.Request) {
		accountId := r.PathValue("accountId")
		accountFindResult, err := accountRepo.Find(DatabaseId(accountId))
		if err != nil {
			panic(err)
		}
		if accountFindResult.IsFound {
			eventId := DatabaseId(r.PathValue("eventId"))
			eventFindResult, err := eventRepo.Find(eventId)
			if err != nil {
				panic(err)
			}
			if eventFindResult.IsFound {
				event := eventFindResult.Record
				event.Name = r.PostFormValue("event[name]")
				event.Date = r.PostFormValue("event[date]")
				event.Hour = r.PostFormValue("event[hour]")
				event.Venue = r.PostFormValue("event[venue]")
				event.Town = r.PostFormValue("event[town]")
				event.Location = r.PostFormValue("event[location]")
				event.Description = r.PostFormValue("event[description]")
				_, err = eventRepo.Save(event)

				if err != nil {
					panic(err)
				}

				http.Redirect(w, r, "/accounts/"+accountId+"/events", http.StatusSeeOther)
			}
		}
	})

	app.router.HandleFunc("DELETE /accounts/{accountId}/events/{eventId}", func(w http.ResponseWriter, r *http.Request) {
		//accountFindResult, err := accountRepo.Find(DatabaseId(r.PathValue("accountId")))
		accountId := r.PathValue("accountId")
		eventId := DatabaseId(r.PathValue("eventId"))
		eventRepo.Delete(eventId)

		http.Redirect(w, r, "/accounts/"+accountId+"/events", http.StatusOK)
	})

	return &controller
}
