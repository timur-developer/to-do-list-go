package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"os"
	"to-do-list-go/internal/database"
	"to-do-list-go/internal/models"
	"to-do-list-go/internal/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=" + os.Getenv("DB_SSLMODE")

	db := database.InitDB(dsn)

	database.MigrateDB(db, &models.Task{})

	e := echo.New()

	router.RegisterRoutes(e, db)

	log.Fatal(e.Start(":8080"))
}
