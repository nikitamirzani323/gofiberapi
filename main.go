package main

import (
	"log"

	"github.com/nikitamirzani323/gofiberapi/db"
	"github.com/nikitamirzani323/gofiberapi/routes"
)

func main() {
	db.Init()
	app := routes.Init()
	log.Fatal(app.Listen(":7071"))
}
