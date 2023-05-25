package main

import (
	"github.com/adYushinW/TestTask/internal/app"
	"github.com/adYushinW/TestTask/internal/db"
	"github.com/adYushinW/TestTask/internal/transport/gin"
)

func main() {

	db, err := db.New()
	if err != nil {
		panic(err)
	}

	app := app.New(db)

	if err := gin.Service(app); err != nil {
		panic(err)
	}
}
