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

	if err := db.AutoMigrate(&models.Project{}, &models.Task{},&models.User{}); err != nil {
		logger.Error("ошибка при выполнении автомиграции", "error", err)
		panic(fmt.Sprintf("не удалось выполнит миграции:%v", err))
	}

	projectRepo := repository.NewProjectRepository(db, logger)
	taskRepo := repository.NewTaskRepository(db, logger)
	userRepo := repository.NewUserRepository(db,logger)

	projectService := service.NewProjectService(db, logger, projectRepo)
	taskService := service.NewTaskService(db, logger, taskRepo)
	userService := service.NewUserService(userRepo,db,logger)

	userHandler := transport.NewUserHandler(userService,logger)

	fmt.Println(projectService, taskService)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	users := r.Group("/users")
	{
		users.POST("/", userHandler.CreateUser) 
	}

	authUsers := r.Group("/users")
	authUsers.Use(middleware.AuthMiddleware(&userRepo))
	{
		authUsers.GET("/:id", userHandler.GetUserByID)
		authUsers.PATCH("/:id", userHandler.UpdateUser)
	}

	adminUsers := r.Group("/admin/users")
	adminUsers.Use(middleware.AuthMiddleware(&userRepo), middleware.RequireAdmin())
	{
		adminUsers.DELETE("/:id", userHandler.DeleteUser)
	}

	logger.Info("Server running on :8080")

	r.Run()
}
