package main

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"github.com/mradmacher/audiofeeler/internal"
)

func main() {
	err := godotenv.Load()
	if err != nil {
      panic("Can't load .env file")
	}
	app, err := audiofeeler.NewApp(audiofeeler.NewTemplateEngine("views"), os.Getenv("AUDIOFEELER_DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer app.Cleanup()

	fmt.Println("Starting the server on :3000...")
	app.Start()
}
