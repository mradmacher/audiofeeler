package main

import (
	"fmt"
	"os"
	"github.com/mradmacher/audiofeeler/internal/web"
)

func main() {
	app, err := web.NewApp("views", os.Getenv("AUDIOFEELER_DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer app.Cleanup()

	fmt.Println("Starting the server on :3000...")
	app.Start()
}
