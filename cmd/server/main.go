package main

import (
	"github.com/labstack/echo/v4"
	"go-rest-api/internal/database"
	"go-rest-api/internal/models"
	"go-rest-api/internal/router"
	"log"
)

/*type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}*/

func main() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable"
	db := database.InitDB(dsn)

	database.MigrateDB(db, &models.Message{})

	e := echo.New()

	router.RegisterRoutes(e, db)

	log.Fatal(e.Start(":8080"))
}
