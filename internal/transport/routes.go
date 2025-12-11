package transport

import (
	"back-minijira-petproject1/internal/service"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	logger *slog.Logger,
	taskService service.TaskService,
	projectService service.ProjectService,
) {
	taskHandler := NewTaskHandler(taskService, logger)
	projectHandler := NewProjectHandler(projectService, logger)
	

	taskHandler.RegisterRoutes(router)
	projectHandler.RegisterRoutes(router)
}
