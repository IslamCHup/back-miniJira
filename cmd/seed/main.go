package main

import (
	"back-minijira-petproject1/internal/config"
	"back-minijira-petproject1/internal/logging"
	"back-minijira-petproject1/internal/models"
	"fmt"
	"log"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	logger := logging.InitLogger()
	db := config.SetUpDatabaseConnection(logger)

	logger.Info("Starting seed data generation...")

	// Очищаем существующие данные (опционально)
	// db.Exec("TRUNCATE TABLE users, projects, tasks, chat_messages CASCADE")

	// Создаем пользователей
	users := createUsers(db, logger)
	if len(users) == 0 {
		log.Fatal("Failed to create users")
	}

	// Создаем проекты
	projects := createProjects(db, logger)
	if len(projects) == 0 {
		log.Fatal("Failed to create projects")
	}

	// Создаем задачи
	tasks := createTasks(db, logger, projects, users)
	if len(tasks) == 0 {
		log.Fatal("Failed to create tasks")
	}

	// Привязываем пользователей к задачам
	assignUsersToTasks(db, logger, tasks, users)

	// Создаем сообщения в чате
	createChatMessages(db, logger, tasks, users)

	logger.Info("Seed data generation completed successfully!")
	fmt.Println("\n✅ Seed данные успешно созданы!")
	fmt.Println("\nТестовые пользователи:")
	fmt.Println("  Админ: admin@test.com / password123")
	fmt.Println("  Пользователь 1: user1@test.com / password123")
	fmt.Println("  Пользователь 2: user2@test.com / password123")
	fmt.Println("  Пользователь 3: user3@test.com / password123")
}

func createUsers(db *gorm.DB, logger *slog.Logger) []models.User {
	users := []models.User{
		{
			FullName:     "Администратор Системы",
			Email:        "admin@test.com",
			PasswordHash: hashPassword("password123"),
			IsAdmin:      true,
			IsVerified:   true,
			VerifyToken:  "",
		},
		{
			FullName:     "Иван Петров",
			Email:        "user1@test.com",
			PasswordHash: hashPassword("password123"),
			IsAdmin:      false,
			IsVerified:   true,
			VerifyToken:  "",
		},
		{
			FullName:     "Мария Сидорова",
			Email:        "user2@test.com",
			PasswordHash: hashPassword("password123"),
			IsAdmin:      false,
			IsVerified:   true,
			VerifyToken:  "",
		},
		{
			FullName:     "Алексей Козлов",
			Email:        "user3@test.com",
			PasswordHash: hashPassword("password123"),
			IsAdmin:      false,
			IsVerified:   true,
			VerifyToken:  "",
		},
	}

	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			logger.Error("Failed to create user", "email", users[i].Email, "error", err)
			continue
		}
		logger.Info("User created", "id", users[i].ID, "email", users[i].Email, "is_admin", users[i].IsAdmin)
	}

	return users
}

func createProjects(db *gorm.DB, logger *slog.Logger) []models.Project {
	now := time.Now()
	futureDate := now.AddDate(0, 3, 0) // +3 месяца

	projects := []models.Project{
		{
			Title:       "Разработка веб-приложения",
			Description: "Создание современного веб-приложения с использованием React и Go. Включает разработку фронтенда, бэкенда и интеграцию с базой данных.",
			Status:      "active",
			TimeEnd:     &futureDate,
		},
		{
			Title:       "Мобильное приложение",
			Description: "Разработка мобильного приложения для iOS и Android. Включает дизайн, разработку и тестирование.",
			Status:      "active",
			TimeEnd:     &futureDate,
		},
		{
			Title:       "Система аналитики",
			Description: "Создание системы аналитики для отслеживания метрик и генерации отчетов. Интеграция с различными источниками данных.",
			Status:      "inactive",
			TimeEnd:     nil,
		},
		{
			Title:       "Завершенный проект",
			Description: "Пример завершенного проекта для тестирования отчетов и статистики.",
			Status:      "completed",
			TimeEnd:     &now,
		},
	}

	for i := range projects {
		if err := db.Create(&projects[i]).Error; err != nil {
			logger.Error("Failed to create project", "title", projects[i].Title, "error", err)
			continue
		}
		logger.Info("Project created", "id", projects[i].ID, "title", projects[i].Title)
	}

	return projects
}

func createTasks(db *gorm.DB, logger *slog.Logger, projects []models.Project, users []models.User) []models.Task {
	if len(projects) == 0 || len(users) == 0 {
		return []models.Task{}
	}

	now := time.Now()
	startDate1 := now.AddDate(0, 0, -10) // 10 дней назад
	startDate2 := now.AddDate(0, 0, -5)  // 5 дней назад
	finishDate1 := now.AddDate(0, 0, 5)   // через 5 дней
	finishDate2 := now.AddDate(0, 0, 10)  // через 10 дней

	tasks := []models.Task{
		// Проект 1: Веб-приложение
		{
			Title:       "Настройка окружения разработки",
			Description: "Установка и настройка всех необходимых инструментов для разработки: IDE, Git, Docker, база данных.",
			Status:      "done",
			ProjectID:   projects[0].ID,
			Priority:    5,
			StartTask:   &startDate1,
			FinishTask:  &now,
		},
		{
			Title:       "Разработка API",
			Description: "Создание REST API для взаимодействия фронтенда и бэкенда. Реализация всех необходимых эндпоинтов.",
			Status:      "in_progress",
			ProjectID:   projects[0].ID,
			Priority:    8,
			StartTask:   &startDate2,
			FinishTask:  &finishDate1,
		},
		{
			Title:       "Разработка UI компонентов",
			Description: "Создание переиспользуемых UI компонентов для фронтенда. Реализация дизайн-системы.",
			Status:      "in_progress",
			ProjectID:   projects[0].ID,
			Priority:    7,
			StartTask:   &startDate2,
			FinishTask:  &finishDate1,
		},
		{
			Title:       "Интеграционное тестирование",
			Description: "Написание и выполнение интеграционных тестов для проверки взаимодействия компонентов системы.",
			Status:      "todo",
			ProjectID:   projects[0].ID,
			Priority:    6,
			StartTask:   nil,
			FinishTask:  &finishDate2,
		},
		// Проект 2: Мобильное приложение
		{
			Title:       "Дизайн макетов",
			Description: "Создание дизайн-макетов для всех экранов мобильного приложения. Проработка UX/UI.",
			Status:      "done",
			ProjectID:   projects[1].ID,
			Priority:    9,
			StartTask:   &startDate1,
			FinishTask:  &now,
		},
		{
			Title:       "Разработка iOS версии",
			Description: "Реализация iOS версии приложения на Swift. Интеграция с бэкенд API.",
			Status:      "in_progress",
			ProjectID:   projects[1].ID,
			Priority:    8,
			StartTask:   &startDate2,
			FinishTask:  &finishDate1,
		},
		{
			Title:       "Разработка Android версии",
			Description: "Реализация Android версии приложения на Kotlin. Синхронизация функционала с iOS версией.",
			Status:      "todo",
			ProjectID:   projects[1].ID,
			Priority:    8,
			StartTask:   nil,
			FinishTask:  &finishDate2,
		},
		// Проект 3: Система аналитики
		{
			Title:       "Проектирование архитектуры",
			Description: "Разработка архитектуры системы аналитики. Выбор технологий и инструментов.",
			Status:      "done",
			ProjectID:   projects[2].ID,
			Priority:    7,
			StartTask:   &startDate1,
			FinishTask:  &now,
		},
		{
			Title:       "Настройка ETL процессов",
			Description: "Настройка процессов извлечения, трансформации и загрузки данных из различных источников.",
			Status:      "todo",
			ProjectID:   projects[2].ID,
			Priority:    6,
			StartTask:   nil,
			FinishTask:  &finishDate2,
		},
		// Проект 4: Завершенный проект
		{
			Title:       "Задача 1 завершенного проекта",
			Description: "Пример выполненной задачи из завершенного проекта.",
			Status:      "done",
			ProjectID:   projects[3].ID,
			Priority:    5,
			StartTask:   &startDate1,
			FinishTask:  &now,
		},
		{
			Title:       "Задача 2 завершенного проекта",
			Description: "Еще одна выполненная задача для тестирования отчетов.",
			Status:      "done",
			ProjectID:   projects[3].ID,
			Priority:    4,
			StartTask:   &startDate1,
			FinishTask:  &now,
		},
	}

	for i := range tasks {
		if err := db.Create(&tasks[i]).Error; err != nil {
			logger.Error("Failed to create task", "title", tasks[i].Title, "error", err)
			continue
		}
		logger.Info("Task created", "id", tasks[i].ID, "title", tasks[i].Title, "project_id", tasks[i].ProjectID)
	}

	return tasks
}

func assignUsersToTasks(db *gorm.DB, logger *slog.Logger, tasks []models.Task, users []models.User) {
	if len(users) < 3 {
		return
	}

	// Привязываем пользователей к задачам
	assignments := []struct {
		taskIndex int
		userIndex int
	}{
		{0, 1}, // Задача 0 -> Пользователь 1
		{1, 1}, // Задача 1 -> Пользователь 1
		{1, 2}, // Задача 1 -> Пользователь 2 (несколько пользователей на задачу)
		{2, 2}, // Задача 2 -> Пользователь 2
		{3, 3}, // Задача 3 -> Пользователь 3
		{4, 1}, // Задача 4 -> Пользователь 1
		{5, 2}, // Задача 5 -> Пользователь 2
		{6, 3}, // Задача 6 -> Пользователь 3
		{7, 1}, // Задача 7 -> Пользователь 1
		{8, 2}, // Задача 8 -> Пользователь 2
		{9, 1}, // Задача 9 -> Пользователь 1
		{9, 2}, // Задача 9 -> Пользователь 2
	}

	for _, assignment := range assignments {
		if assignment.taskIndex >= len(tasks) || assignment.userIndex >= len(users) {
			continue
		}

		task := tasks[assignment.taskIndex]
		user := users[assignment.userIndex]

		if err := db.Model(&task).Association("Users").Append(&user); err != nil {
			logger.Error("Failed to assign user to task", "task_id", task.ID, "user_id", user.ID, "error", err)
			continue
		}
		logger.Info("User assigned to task", "task_id", task.ID, "user_id", user.ID)
	}
}

func createChatMessages(db *gorm.DB, logger *slog.Logger, tasks []models.Task, users []models.User) {
	if len(tasks) == 0 || len(users) < 2 {
		return
	}

	messages := []models.ChatMessage{
		// Сообщения для задачи 1 (in_progress)
		{
			UserID:      users[1].ID,
			Text:        "Начал работу над API. Планирую реализовать основные эндпоинты к концу недели.",
			ChatableID:  tasks[1].ID,
			ChatableType: "tasks",
		},
		{
			UserID:      users[2].ID,
			Text:        "Отлично! Если нужна помощь с тестированием, дай знать.",
			ChatableID:  tasks[1].ID,
			ChatableType: "tasks",
		},
		{
			UserID:      users[1].ID,
			Text:        "Спасибо! Обязательно обращусь.",
			ChatableID:  tasks[1].ID,
			ChatableType: "tasks",
		},
		// Сообщения для задачи 2
		{
			UserID:      users[2].ID,
			Text:        "Создал базовые компоненты: кнопки, формы, модальные окна. Что дальше?",
			ChatableID:  tasks[2].ID,
			ChatableType: "tasks",
		},
		{
			UserID:      users[1].ID,
			Text:        "Можешь перейти к компонентам списков и таблиц.",
			ChatableID:  tasks[2].ID,
			ChatableType: "tasks",
		},
		// Сообщения для задачи 5
		{
			UserID:      users[2].ID,
			Text:        "iOS версия почти готова. Осталось протестировать на реальных устройствах.",
			ChatableID:  tasks[5].ID,
			ChatableType: "tasks",
		},
		{
			UserID:      users[0].ID, // Админ
			Text:        "Хорошая работа! После тестирования можно будет переходить к Android версии.",
			ChatableID:  tasks[5].ID,
			ChatableType: "tasks",
		},
	}

	for i := range messages {
		if err := db.Create(&messages[i]).Error; err != nil {
			logger.Error("Failed to create chat message", "task_id", messages[i].ChatableID, "error", err)
			continue
		}
		logger.Info("Chat message created", "id", messages[i].ID, "task_id", messages[i].ChatableID)
	}
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}
	return string(hash)
}

