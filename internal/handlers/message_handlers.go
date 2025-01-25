package handlers

import (
	"github.com/go-playground/validator/v10"
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

	return respondWithSuccess(c, http.StatusOK, messages)
}

var validate = validator.New()

func PostHandler(c echo.Context, db *gorm.DB) error {
	var message models.Message
	if err := c.Bind(&message); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Could not add the message")
	}

	if err := validate.Struct(&message); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Validation failed")
	}

	if err := db.Create(&message).Error; err != nil {
		return respondWithError(c, http.StatusBadRequest, "Could not create the message")
	}

	return respondWithSuccess(c, http.StatusOK, createResponse("OK", "Message was added successfully"))
}

func PutHandler(c echo.Context, db *gorm.DB) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, "Incorrect ID")
	}
	var message models.Message
	if err := db.First(&message, id).Error; err != nil {
		return respondWithError(c, http.StatusNotFound, "Could not find the message")
	}
	var updatedMessage models.Message
	if err := c.Bind(&updatedMessage); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid input")
	}

	if err := db.Model(&message).Updates(&updatedMessage).Error; err != nil {
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

	var message models.Message
	if err := db.First(&message, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondWithError(c, http.StatusNotFound, "Message not found")
		}
		return respondWithError(c, http.StatusBadRequest, "Could not retrieve the message")
	}

	if err := db.Delete(&message).Error; err != nil {
		return respondWithError(c, http.StatusBadRequest, "Could not delete the message")
	}

	return respondWithSuccess(c, http.StatusOK, createResponse("Success", "Message was deleted"))
}
