package main

import (
	"fmt"
	"github.com/mradmacher/audiofeeler/internal"
)

func main() {
	app, err := audiofeeler.NewApp("views")
	if err != nil {
		panic(err)
	}
	defer app.Cleanup()

	fmt.Println("Starting the server on :3000...")
	app.Start()
}
