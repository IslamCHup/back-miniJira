package transport

import (
	"back-minijira-petproject1/internal/middleware"
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

func (h *TaskHandler) RegisterRoutes(r *gin.Engine, authService service.AuthService) {
	authTasks := r.Group("/tasks")
	authTasks.Use(middleware.AuthMiddleware(authService))
	{
		authTasks.GET("/", h.ListTasks)
		authTasks.GET("/:id", h.GetTaskByID)
		authTasks.POST("/:id/assign", h.AssignTask)
		authTasks.POST("/:id/unassign", h.UnassignTask)
	}
	adminTasks := r.Group("/admin/tasks")
	adminTasks.Use(middleware.AuthMiddleware(authService), middleware.RequireAdmin())
	{
		adminTasks.POST("/", h.Create)
		adminTasks.PATCH("/:id", h.Update)
		adminTasks.DELETE("/:id", h.DeleteTask)
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
	filter := models.TaskFilter{
		Limit:  20,
		Offset: 0,
	}

	// Обрабатываем параметры запроса только если они не пустые
	if status := strings.ToLower(strings.TrimSpace(c.Query("status"))); status != "" {
		filter.Status = &status
	}

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userIDInt, err := strconv.Atoi(userIDStr); err == nil && userIDInt > 0 {
			userID := uint(userIDInt)
			filter.UserID = &userID
		}
	}

	if projectIDStr := c.Query("project_id"); projectIDStr != "" {
		if projectIDInt, err := strconv.Atoi(projectIDStr); err == nil && projectIDInt > 0 {
			projectID := uint(projectIDInt)
			filter.ProjectID = &projectID
		}
	}

	if search := c.Query("search"); search != "" {
		filter.Search = &search
	}

	if priorityStr := c.Query("priority"); priorityStr != "" {
		if priority, err := strconv.Atoi(priorityStr); err == nil {
			filter.Priority = &priority
		}
	}

	if sortBy := strings.ToLower(strings.TrimSpace(c.Query("sort_by"))); sortBy != "" {
		filter.SortBy = &sortBy
	}

	if sortOrder := strings.ToLower(strings.TrimSpace(c.Query("sort_order"))); sortOrder != "" {
		filter.SortOrder = &sortOrder
	}

	tasks, err := h.service.ListTasks(&filter)

	if err != nil {
		h.logger.Error("failed to list tasks", "op", "task.handler.ListTasks", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tasks: " + err.Error()})
		return
	}

	h.logger.Info("tasks listed", "op", "task.handler.ListTasks", "count", len(tasks))
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

	// Поддержка camelCase поля от фронта
	if req.ProjectID == 0 && req.ProjectId != 0 {
		req.ProjectID = req.ProjectId
	}

	if req.ProjectID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
		return
	}

	// Значение по умолчанию для статуса
	if strings.TrimSpace(req.Status) == "" {
		req.Status = "todo"
	}

	if err := h.service.CreateTask(&req); err != nil {
		h.logger.Error("failed to create task", "op", "task.handler.Create", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}

	h.logger.Info("task created", "op", "task.handler.Create", "project_id", req.ProjectID)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update successful"})
}

func (h *TaskHandler) AssignTask(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.Param("id"))
	currentUser := c.MustGet("currentUser").(models.User)

	if err := h.service.AssignTaskToUser(uint(taskID), currentUser.ID); err != nil {
		h.logger.Error("AssignTask failed", "task_id", taskID, "user_id", currentUser.ID, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("AssignTask success", "task_id", taskID, "user_id", currentUser.ID)
	c.JSON(http.StatusOK, gin.H{"message": "task assigned successfully"})
}

func (h *TaskHandler) UnassignTask(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.Param("id"))
	currentUser := c.MustGet("currentUser").(models.User)

	if err := h.service.UnassignTaskFromUser(uint(taskID), currentUser.ID); err != nil {
		h.logger.Error("UnassignTask failed", "task_id", taskID, "user_id", currentUser.ID, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("UnassignTask success", "task_id", taskID, "user_id", currentUser.ID)
	c.JSON(http.StatusOK, gin.H{"message": "task unassigned successfully"})
}
