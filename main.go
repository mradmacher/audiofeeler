package main

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"github.com/mradmacher/audiofeeler/internal"
)

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
      panic("Can't load .env file")
	}
	dbEngine, err := internal.NewDbEngine(os.Getenv("AUDIOFEELER_DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	templateEngine := internal.NewTemplateEngine("views")
	app := internal.NewApp(templateEngine, dbEngine)
	defer app.Cleanup()

	fmt.Println("Starting the server on :3000...")
	app.Start()
}
