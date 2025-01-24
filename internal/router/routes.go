package router

import (
	"github.com/labstack/echo/v4"
	"go-rest-api/internal/handlers"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	e.GET("/messages", func(c echo.Context) error {
		return handlers.GetHandler(c, db)
	})
	e.POST("/messages", func(c echo.Context) error {
		return handlers.PostHandler(c, db)
	})
	e.PATCH("/messages/:id", func(c echo.Context) error {
		return handlers.PatchHandler(c, db)
	})
	e.DELETE("/messages/:id", func(c echo.Context) error {
		return handlers.DeleteHandler(c, db)
	})
}
