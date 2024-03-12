package main

import (
	"fmt"
	"github.com/mradmacher/audiofeeler/internal"
	"os"
)

func main() {
	db, err := audiofeeler.NewDbClient(os.Getenv("AUDIOFEELER_DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//defer func() {
	//    db.RemoveStructure()
	//    fmt.Println("Tables dropped")
	//}()

	r := audiofeeler.EventsRepo{db}
	events, err := r.FindAll()
	if err != nil {
		panic(err)
	}

	for _, event := range *events {
		fmt.Printf("Event[%v]: %v\n", event.Date.Value(), event)
	}
}
