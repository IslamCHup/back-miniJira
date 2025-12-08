package main

import (
	"back-minijira-petproject1/internal/config"
	"back-minijira-petproject1/internal/logging"
	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/repository"
	"back-minijira-petproject1/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := logging.InitLogger()

	db := config.SetUpDatabaseConnection(logger)

	if err := db.AutoMigrate(&models.Project{}, &models.Task{}); err != nil {
		logger.Error("ошибка при выполнении автомиграции", "error", err)
		panic(fmt.Sprintf("не удалось выполнит миграции:%v", err))
	}

	projectRepo := repository.NewProjectRepository(db, logger)
	taskRepo := repository.NewTaskRepository(db, logger)

	projectService := service.NewProjectService(db, logger, projectRepo)
	taskService := service.NewTaskService(db, logger, taskRepo)

	fmt.Println(projectService, taskService)

	r := gin.Default()

	r.Run()
}
