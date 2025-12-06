package main

import (
	"back-minijira-petproject1/internal/config"
	"back-minijira-petproject1/internal/logging"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := logging.InitLogger()
	
	db := config.SetUpDatabaseConnection(logger)

	if err := db.AutoMigrate(); err != nil {
		logger.Error("ошибка при выполнении автомиграции","error",err)
		panic(fmt.Sprintf("не удалось выполнит миграции:%v",err))
	}

r := gin.Default()


r.Run()
}