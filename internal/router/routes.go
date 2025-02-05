package router

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"to-do-list-go/internal/handlers"
	"to-do-list-go/internal/kafka"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, producer *kafka.Producer) {
	e.GET("/list", func(c echo.Context) error {
		return handlers.GetTasksHandler(c, db, producer)
	})
	e.POST("/create", func(c echo.Context) error {
		return handlers.PostTasksHandler(c, db, producer)
	})
	e.PATCH("/done/:id", func(c echo.Context) error { return handlers.PatchTasksHandler(c, db, producer) })
	e.DELETE("/delete/:id", func(c echo.Context) error {
		return handlers.DeleteTasksHandler(c, db, producer)
	})
}
