package handlers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"to-do-list-go/internal/models"
)

func respondWithError(c echo.Context, status int, task string) error {
	return c.JSON(status, map[string]string{"error": task})
}

func respondWithSuccess(c echo.Context, status int, data interface{}) error {
	return c.JSON(status, data)
}

func createResponse(status string, task string) map[string]string {
	return map[string]string{
		"status": status,
		"task":   task,
	}
}

func GetTasksHandler(c echo.Context, db *gorm.DB) error {
	var tasks []models.Task
	if err := db.Find(&tasks).Error; err != nil {
		return respondWithError(c, http.StatusInternalServerError, "Error retrieving task")
	}
	if len(tasks) == 0 {
		return respondWithSuccess(c, http.StatusOK, map[string]string{"task": "No task found"})
	}
	return respondWithSuccess(c, http.StatusOK, tasks)
}

func PostTasksHandler(c echo.Context, db *gorm.DB) error {
	var task models.Task
	if err := c.Bind(&task); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Could not add the task")
	}

	if err := db.Create(&task).Error; err != nil {
		return respondWithError(c, http.StatusBadRequest, "Could not create the task")
	}

	return respondWithSuccess(c, http.StatusOK, createResponse("OK", "Task was added successfully"))
}

func PutTasksHandler(c echo.Context, db *gorm.DB) error {
	// Получение ID из параметров URL
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, "Incorrect ID format")
	}

	// Проверка наличия записи в базе
	var task models.Task
	if err := db.First(&task, id).Error; err != nil {
		return respondWithError(c, http.StatusNotFound, "task not found")
	}

	// Получение данных для обновления
	var updatedTask models.Task
	if err := c.Bind(&updatedTask); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid input format")
	}

	// Обновление только изменённых полей
	updates := map[string]interface{}{}
	if updatedTask.TaskName != "" {
		updates["task_name"] = updatedTask.TaskName
	}
	if updatedTask.TaskDescription != "" {
		updates["description"] = updatedTask.TaskDescription
	}
	updates["is_done"] = updatedTask.IsDone

	updates["status_updated_at"] = gorm.Expr("CURRENT_TIMESTAMP")

	// Обновление записи в базе данных
	if err := db.Model(&task).Where("id = ?", id).Updates(updates).Error; err != nil {
		return respondWithError(c, http.StatusInternalServerError, "Could not update the task")
	}

	// Возвращение успешного ответа
	return respondWithSuccess(c, http.StatusOK, createResponse("Success", "Task was updated"))
}

func DeleteTasksHandler(c echo.Context, db *gorm.DB) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, "Incorrect ID")
	}

	// Проверяем, существует ли сообщение с таким ID
	var task models.Task
	if err := db.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondWithError(c, http.StatusNotFound, "Task not found")
		}
		return respondWithError(c, http.StatusBadRequest, "Could not retrieve the task")
	}

	// Если сообщение найдено, удаляем его
	if err := db.Delete(&task).Error; err != nil {
		return respondWithError(c, http.StatusBadRequest, "Could not delete the task")
	}

	return respondWithSuccess(c, http.StatusOK, createResponse("Success", "Task was deleted"))
}
