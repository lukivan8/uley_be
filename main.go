package main

import (
	"log"

	_ "uley_be/migrations"
	appRouter "uley_be/router"

	"github.com/pocketbase/pocketbase"
)

func main() {
	app := pocketbase.New()

	appRouter.RegisterRoutes(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
