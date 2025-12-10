package transport

import (
	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/service"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service service.TaskService
	logger  *slog.Logger
}

func NewTaskHandler(service service.TaskService, logger *slog.Logger) TaskHandler {
	return TaskHandler{service: service, logger: logger}
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	task, err := h.service.GetTaskByID(uint(id))

	if err != nil {
		h.logger.Error("failed to get task by id", "op", "task.handler.GetByID", "id", id, "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get task by id"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) ListTasks(c *gin.Context) {

	// --- status ---
	status := c.Query("status")

	// --- user_id ---
	var userID *uint
	if v := c.Query("user_id"); v != "" {
		if idInt, _ := strconv.Atoi(v) {
			tmp := uint(idInt)
			userID = &tmp
		}
	}

	// --- project_id ---
	var projectID *uint
	if v := c.Query("project_id"); v != "" {
		if idInt, err := strconv.Atoi(v); err == nil {
			tmp := uint(idInt)
			projectID = &tmp
		}
	}

	// --- search ---
	var searchPtr *string
	if v := c.Query("search"); v != "" {
		searchPtr = &v
	}

	// --- priority (int) ---
	var priorityPtr *int
	if v := c.Query("priority"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			priorityPtr = &p
		}
	}

	// --- sort ---
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")

	// --- pagination ---
	var limitPtr *int
	if v := c.Query("limit"); v != "" {
		if x, err := strconv.Atoi(v); err == nil {
			limitPtr = &x
		}
	}

	var offsetPtr *int
	if v := c.Query("offset"); v != "" {
		if x, err := strconv.Atoi(v); err == nil {
			offsetPtr = &x
		}
	}

	// --- собрать фильтр ---
	filter := models.TaskFilter{
		Status:    &status,
		UserID:    userID,
		ProjectID: projectID,
		Search:    searchPtr,
		Priority:  priorityPtr,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Limit:     limitPtr,
		Offset:    offsetPtr,
	}

	// --- применить фильтр ---
	db := h.db
	db = ApplyTaskFilter(db, filter)

	// пагинация
	if limitPtr != nil {
		db = db.Limit(*limitPtr)
	}
	if offsetPtr != nil {
		db = db.Offset(*offsetPtr)
	}

	// выполнить запрос
	var tasks []models.Task
	if err := db.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// конвертация в response (priority → string)
	resp := make([]TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		resp = append(resp, ConvertTaskToResponse(t))
	}

	c.JSON(http.StatusOK, resp)
}

/*
	ListTasks(filter *models.TaskFilter) ([]*models.TaskResponse, error)
	DeleteTask(id uint) error
	CreateTask(req models.TaskCreateReq) error
	UpdateTask(id uint, req models.TaskUpdateReq) error
*/
