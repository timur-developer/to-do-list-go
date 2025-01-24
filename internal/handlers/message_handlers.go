package handlers

import (
	"github.com/labstack/echo/v4"
	"go-rest-api/internal/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func respondWithError(c echo.Context, status int, message string) error {
	return c.JSON(status, map[string]string{"error": message})
}

func respondWithSuccess(c echo.Context, status int, data interface{}) error {
	return c.JSON(status, data)
}

func createResponse(status string, message string) map[string]string {
	return map[string]string{
		"status":  status,
		"message": message,
	}
}

func GetHandler(c echo.Context, db *gorm.DB) error {
	var messages []models.Message

	if err := db.Find(&messages).Error; err != nil {
		return respondWithError(c, http.StatusBadRequest, "Could not find the message")
	}

	return respondWithSuccess(c, http.StatusOK, &messages)
}

func PostHandler(c echo.Context, db *gorm.DB) error {
	var message models.Message
	if err := c.Bind(&message); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Could not add the message")
	}

	if err := db.Create(&message).Error; err != nil {
		return respondWithError(c, http.StatusBadRequest, "Could not create the message")
	}

	return respondWithSuccess(c, http.StatusOK, createResponse("OK", "Message was added successfully"))
}

func PatchHandler(c echo.Context, db *gorm.DB) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, "Incorrect ID")
	}
	var updatedMessage models.Message
	if err := c.Bind(&updatedMessage); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid input")
	}

	if err := db.Model(&models.Message{}).Where("id = ?", id).Update("text", updatedMessage.Text).Error; err != nil {
		return respondWithError(c, http.StatusBadRequest, "Could not update the message")
	}

	return respondWithSuccess(c, http.StatusOK, createResponse("Success", "Message was updated"))
}

func DeleteHandler(c echo.Context, db *gorm.DB) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, "Incorrect ID")
	}

	if err := db.Delete(&models.Message{}, id).Error; err != nil {
		return respondWithError(c, http.StatusBadRequest, "Could not delete the message")
	}

	return respondWithSuccess(c, http.StatusOK, createResponse("Success", "Message was deleted"))
}
