package audiofeeler

import (
	"html/template"
)

type EventPresenter struct {
	Account Account
	Event   Event
}

type EventsPresenter struct {
	Account Account
	Events  []Event
}

type AccountsPresenter struct {
	Accounts []Account
}

type EventView struct {
	indexTemplate *template.Template
	newTemplate   *template.Template
}

type AccountView struct {
	indexTemplate *template.Template
	showTemplate  *template.Template
}

func NewAccountView(t TemplateEngine) *AccountView {
	view := AccountView{}
	view.indexTemplate = t.Parse("accounts")
	view.showTemplate = t.Parse("account", "account_wrapper")
	return &view
}

func (view *AccountView) renderShow(vc ViewContext, account Account) error {
	return vc.execute(view.showTemplate, account)
}

func (view *AccountView) renderIndex(vc ViewContext, accounts []Account) error {
	return vc.execute(view.indexTemplate, AccountsPresenter{accounts})
}

func NewEventView(t TemplateEngine) *EventView {
	view := EventView{}

	view.indexTemplate = t.Parse("events", "account_wrapper")
	view.newTemplate = t.Parse("event_form", "account_wrapper")

	return &view
}

func (view *EventView) renderIndex(vc ViewContext, account Account, events []Event) error {
	return vc.execute(view.indexTemplate, EventsPresenter{account, events})
}

func (view *EventView) renderNew(vc ViewContext, account Account, event Event) error {
	return vc.execute(view.newTemplate, EventPresenter{account, event})
}
