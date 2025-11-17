package main

import (
	"bytes"
	"fmt"
	"github.com/mradmacher/audiofeeler/internal"
	"github.com/mradmacher/audiofeeler/optiomist"
	"os"
)

func exampleAccounts() []audiofeeler.Account {
	return []audiofeeler.Account {
		audiofeeler.Account {
			Title: optiomist.Some("Czarny Motyl"),
			Name: optiomist.Some("czarnymotyl"),
			Url: optiomist.Some("http://czarnymotyl.art"),
		},
		audiofeeler.Account {
			Title: optiomist.Some("Karoryfer Lecolds"),
			Name: optiomist.Some("karoryfer"),
			Url: optiomist.Some("http://karoryfer.com"),
		},
		audiofeeler.Account {
			Title: optiomist.Some("BalkanArtz"),
			Name: optiomist.Some("balkanartz"),
			Url: optiomist.Some("http://balkanartz.eu"),
		},
		audiofeeler.Account {
			Title: optiomist.Some("Iglika"),
			Name: optiomist.Some("iglika"),
		},
	}
}

func createExampleData(db *audiofeeler.DbClient) {
	r := audiofeeler.AccountsRepo{db}

	for _, account := range exampleAccounts() {
		id, err := r.Create(account)
		fmt.Printf("Account created [%v]: %v\n", id, err)
	}
}

func loadEvents(db *audiofeeler.DbClient, fileName string) {
	jsonBlob, err := os.ReadFile(fileName)
	events, err := audiofeeler.LoadEvents(bytes.NewReader(jsonBlob))
	if err != nil {
		panic(err)
	}

	r := audiofeeler.EventsRepo{db}

	for _, event := range events {
		id, err := r.Create(event)
		fmt.Printf("Event created [%v]: %v\n", id, err)
	}
}

func main() {
	db, err := audiofeeler.NewDbClient(os.Getenv("AUDIOFEELER_DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if db.Ping() {
		fmt.Println("Connected to database")
	} else {
		panic("Not connected to database")
	}

	err = db.RemoveStructure()
	if err != nil {
		panic(err)
	}
	err = db.CreateStructure()
	if err != nil {
		panic(err)
	}
	fmt.Println("Tables created")

	if len(os.Args) < 2 {
		createExampleData(db)
	} else {
		fileName := os.Args[1]
		loadEvents(db, fileName)
	}
}
