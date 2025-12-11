package main

import (
	"back-minijira-petproject1/internal/config"
	"back-minijira-petproject1/internal/logging"
	"back-minijira-petproject1/internal/middleware"
	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/repository"
	"back-minijira-petproject1/internal/service"
	"back-minijira-petproject1/internal/transport"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := logging.InitLogger()

	db := config.SetUpDatabaseConnection(logger)

	// db.Migrator().DropTable(&models.User{})
	if err := db.AutoMigrate(&models.Project{}, &models.Task{}, &models.User{}, &models.ChatMessage{}); err != nil {
		logger.Error("ошибка при выполнении автомиграции", "error", err)
		panic(fmt.Sprintf("не удалось выполнит миграции:%v", err))
	}

	projectRepo := repository.NewProjectRepository(db, logger)
	taskRepo := repository.NewTaskRepository(db, logger)
	userRepo := repository.NewUserRepository(db, logger)
	reportRepo := repository.NewReportRepository(db, logger)
	chatRepo := repository.NewChatRepositoryGorm(db)

	projectService := service.NewProjectService(db, logger, projectRepo)
	taskService := service.NewTaskService(db, logger, taskRepo, projectRepo)
	userService := service.NewUserService(userRepo, db, logger)
	reportService := service.NewReportService(reportRepo, logger)
	chatService := service.NewChatService(chatRepo, logger)
	authService := service.NewAuthService(userRepo, logger)
	r := gin.Default()

	// Добавляем CORS middleware
	r.Use(middleware.CORS())

	transport.RegisterRoutes(
		r, logger, taskService, projectService, reportService, chatService, userService, authService, userRepo,
	)

	logger.Info("Server running on :8080")

	r.Run(":8080")
}
