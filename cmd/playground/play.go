package main

import (
	"fmt"
	"github.com/mradmacher/audiofeeler/internal/repo"
	"os"
)

func main() {
	db, err := repo.Connect(os.Getenv("AUDIOFEELER_DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//defer func() {
	//    db.RemoveStructure()
	//    fmt.Println("Tables dropped")
	//}()

	r := repo.EventsRepo{db}
	events, err := r.All()
	if err != nil {
		panic(err)
	}

	for _, event := range *events {
		fmt.Printf("Event[%v]: %v\n", event.Date.Value(), event)
	}
}
