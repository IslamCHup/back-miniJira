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
	reportService service.ReportService,
	chatService service.ChatService,
) {
	taskHandler := NewTaskHandler(taskService, logger)
	projectHandler := NewProjectHandler(projectService, logger)
	reportHandler := NewReportHandler(reportService, logger)
	chatHandler := NewChatHandler(chatService, logger)

	chatHandler.SetupChatRoutes(router)
	reportHandler.RegisterRoutes(router)
	taskHandler.RegisterRoutes(router)
	projectHandler.RegisterRoutes(router)
}
