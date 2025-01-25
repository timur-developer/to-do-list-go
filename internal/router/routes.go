package router

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"to-do-list-go/internal/handlers"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	e.GET("/list", func(c echo.Context) error {
		return handlers.GetHandler(c, db)
	})
	e.POST("/create", func(c echo.Context) error {
		return handlers.PostHandler(c, db)
	})
	e.PUT("/done/:id", func(c echo.Context) error { return handlers.PutHandler(c, db) })
	e.DELETE("/delete/:id", func(c echo.Context) error {
		return handlers.DeleteHandler(c, db)
	})
}
