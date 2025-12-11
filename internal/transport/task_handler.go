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

func (h *TaskHandler) RegisterRoutes(r *gin.Engine) {
	tasks := r.Group("/tasks")
	{
		tasks.GET("/", h.ListTasks)
		tasks.GET("/:id", h.GetTaskByID)
		tasks.POST("/", h.Create)
		tasks.PATCH("/:id", h.Update)
		tasks.DELETE("/:id", h.DeleteTask)
	}
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
	status := strings.ToLower(strings.TrimSpace(c.Query("status")))

	userIDInt, _ := strconv.Atoi(c.Query("user_id"))
	userID := uint(userIDInt)

	projectIDInt, _ := strconv.Atoi(c.Query("project_id"))
	projectID := uint(projectIDInt)

	search := c.Query("search")

	priority, _ := strconv.Atoi(c.Query("priority"))

	sortBy := strings.ToLower(strings.TrimSpace(c.Query("sort_by")))

	sortOrder := strings.ToLower(strings.TrimSpace(c.Query("sort_Order")))

	filter := models.TaskFilter{
		Status:    &status,
		UserID:    &userID,
		ProjectID: &projectID,
		Search:    &search,
		Priority:  &priority,
		SortBy:    &sortBy,
		SortOrder: &sortOrder,
		Limit:     20,
		Offset:    0,
	}

	tasks, err := h.service.ListTasks(&filter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.DeleteTask(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	h.logger.Info("task deleted", "op", "task.handler.Delete", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "delete successful"})
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req models.TaskCreateReq

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid create task body", "op", "task.handler.Create", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	if err := h.service.CreateTask(&req); err != nil {
		h.logger.Error("failed to create task", "op", "task.handler.Create", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}

	h.logger.Info("task created", "op", "task.handler.Create")
	c.JSON(http.StatusCreated, gin.H{"message": "create successful"})
}

func (h *TaskHandler) Update(c *gin.Context) {
	var req models.TaskUpdateReq

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid update task body", "op", "task.handler.Update", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.UpdateTask(uint(id), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update successful"})
}
