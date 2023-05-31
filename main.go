package main

import (
	"github.com/adYushinW/TestTask/internal/app"
	"github.com/adYushinW/TestTask/internal/db"
	"github.com/adYushinW/TestTask/internal/transport/http"
)

func main() {

	db, err := db.New()
	if err != nil {
		panic(err)
	}

	app := app.New(db)

	if err := http.Service(app); err != nil {
		panic(err)
	}
}
