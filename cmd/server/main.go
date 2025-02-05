package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"os"
	"to-do-list-go/internal/database"
	"to-do-list-go/internal/kafka"
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

	// Kafka
	brokers := []string{"kafka:9092"}
	producer := kafka.NewProducer(brokers)

	// Запуск консьюмера в отдельной горутине
	go kafka.StartConsumer(brokers, "task_updates")

	e := echo.New()
	router.RegisterRoutes(e, db, producer)

	defer producer.Close()

	log.Println("Server is running on port 8080...")
	log.Fatal(e.Start(":8080"))
}
