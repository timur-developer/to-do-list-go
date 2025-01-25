package main

import (
	"github.com/labstack/echo/v4"
	"log"
	"to-do-list-go/internal/database"
	"to-do-list-go/internal/models"
	"to-do-list-go/internal/router"
)

func main() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable"
	db := database.InitDB(dsn)

	database.MigrateDB(db, &models.Message{})

	e := echo.New()

	router.RegisterRoutes(e, db)

	log.Fatal(e.Start(":8080"))
}
