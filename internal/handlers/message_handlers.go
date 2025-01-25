package handlers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"to-do-list-go/internal/models"
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
		return respondWithError(c, http.StatusInternalServerError, "Error retrieving messages")
	}
	if len(messages) == 0 {
		return respondWithSuccess(c, http.StatusOK, map[string]string{"message": "No messages found"})
	}
	return respondWithSuccess(c, http.StatusOK, messages)
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

func PutHandler(c echo.Context, db *gorm.DB) error {
	// Получение ID из параметров URL
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, "Incorrect ID format")
	}

	// Проверка наличия записи в базе
	var message models.Message
	if err := db.First(&message, id).Error; err != nil {
		return respondWithError(c, http.StatusNotFound, "Message not found")
	}

	// Получение данных для обновления
	var updatedMessage models.Message
	if err := c.Bind(&updatedMessage); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid input format")
	}

	// Обновление только изменённых полей
	updates := map[string]interface{}{}
	if updatedMessage.Text != "" {
		updates["text"] = updatedMessage.Text
	}
	updates["is_done"] = updatedMessage.IsDone

	// Обновление записи в базе данных
	if err := db.Model(&message).Where("id = ?", id).Updates(updates).Error; err != nil {
		return respondWithError(c, http.StatusInternalServerError, "Could not update the message")
	}

	// Возвращение успешного ответа
	return respondWithSuccess(c, http.StatusOK, createResponse("Success", "Message was updated"))
}

func DeleteHandler(c echo.Context, db *gorm.DB) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, "Incorrect ID")
	}

	// Проверяем, существует ли сообщение с таким ID
	var message models.Message
	if err := db.First(&message, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondWithError(c, http.StatusNotFound, "Message not found")
		}
		return respondWithError(c, http.StatusBadRequest, "Could not retrieve the message")
	}

	// Если сообщение найдено, удаляем его
	if err := db.Delete(&message).Error; err != nil {
		return respondWithError(c, http.StatusBadRequest, "Could not delete the message")
	}

	return respondWithSuccess(c, http.StatusOK, createResponse("Success", "Message was deleted"))
}
