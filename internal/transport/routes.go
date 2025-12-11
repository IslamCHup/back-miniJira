package transport

import (
	"back-minijira-petproject1/internal/repository"
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
	userService service.UserService,
	authService service.AuthService,
	userRepo repository.UserRepository,
) {
	taskHandler := NewTaskHandler(taskService, logger)
	projectHandler := NewProjectHandler(projectService, logger)
	reportHandler := NewReportHandler(reportService, logger)
	chatHandler := NewChatHandler(chatService, logger)
	userHandler := NewUserHandler(userService, logger)
	authHandler := NewAuthHandler(&authService, logger)
	

	chatHandler.SetupChatRoutes(router)
	reportHandler.RegisterRoutes(router)
	taskHandler.RegisterRoutes(router)
	projectHandler.RegisterRoutes(router)
	userHandler.RegisterRoutes(router, userRepo)
	authHandler.SetupRoutes(router)
	
}
