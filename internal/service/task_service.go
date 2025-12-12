package service

import (
	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/repository"
	"errors"
	"log/slog"
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"
)

type TaskService interface {
	GetTaskByID(id uint) (*models.TaskResponse, error)
	ListTasks(filter *models.TaskFilter) ([]*models.TaskResponse, error)
	DeleteTask(id uint) error
	CreateTask(req *models.TaskCreateReq) error
	UpdateTask(id uint, req models.TaskUpdateReq) error
	AssignTaskToUser(taskID uint, userID uint) error
	UnassignTaskFromUser(taskID uint, userID uint) error
}

type taskService struct {
	db          *gorm.DB
	logger      *slog.Logger
	repo        repository.TaskRepository
	projectRepo repository.ProjectRepository
}

func NewTaskService(db *gorm.DB, logger *slog.Logger, repo repository.TaskRepository, projectRepo repository.ProjectRepository) TaskService {
	return &taskService{db: db, logger: logger, repo: repo, projectRepo: projectRepo}
}

func (s *taskService) GetTaskByID(id uint) (*models.TaskResponse, error) {
	task, err := s.repo.GetTaskByID(id)
	if err != nil {
		s.logger.Error("failed get task by id",
			"op", "service.task.GetTaskByID",
			"id", id,
			"error", err)
		return nil, err
	}

	taskResponse := buildTaskResponse(task)

	s.logger.Info("get task by id successful", "op", "service.task.GetTaskByID", "task", task)
	return taskResponse, nil
}

func (s *taskService) ListTasks(filter *models.TaskFilter) ([]*models.TaskResponse, error) {
	tasks, err := s.repo.ListTasks(filter)

	if err != nil {
		s.logger.Error("empty inquiry",
			"op", "service.task.ListTasks",
			"error", err)
		return nil, err
	}
	tasksResponse := []*models.TaskResponse{}
	for _, task := range tasks {
		taskTransport := buildTaskResponse(task)
		tasksResponse = append(tasksResponse, taskTransport)
	}

	s.logger.Info("list tasks successful", "op", "service.task.ListTasks")
	return tasksResponse, nil
}

func (s *taskService) DeleteTask(id uint) error {
	if err := s.repo.DeleteTask(id); err != nil {
		s.logger.Error("failed delete task by id", "id", id, "err", err)
		return err
	}
	s.logger.Info("delete task by id successful", "op", "service.project.DeleteTask")
	return nil
}

func (s *taskService) CreateTask(req *models.TaskCreateReq) error {
	if err := s.repo.CreateTask(req); err != nil {
		s.logger.Error("failed create task from req", "err", err, "req", req)
		return err
	}
	s.logger.Info("create task from req successful", "op", "service.project.CreateTask")
	return nil
}

func (s *taskService) UpdateTask(id uint, req models.TaskUpdateReq) error {
	s.logger.Info("UpdateTask called", "op", "service.task.UpdateTask", "id", id,
		"title", req.Title, "status", req.Status, "priority", req.Priority)

	return s.db.Transaction(func(tx *gorm.DB) error {
		taskrepo := s.repo.WithDB(tx)
		projectRepo := s.projectRepo.WithDB(tx)

		task, err := taskrepo.GetTaskByID(id)
		if err != nil {
			s.logger.Error("failed to get the task by id",
				"op", "service.task.UpdateTask", "id", id, "err", err)
			return err
		}

		s.logger.Info("task found", "op", "service.task.UpdateTask", "task_id", task.ID, "current_status", task.Status)

		//этого блого нет после ревью
		allowedTransport := map[string][]string{
			"todo":        {"in_progress"},
			"in_progress": {"todo", "done"},
			"done":        {"in_progress"},
		}

		oldStatusTask := strings.ToLower(strings.TrimSpace(task.Status))
		var newStatusTask string
		if req.Status != nil {
			newStatusTask = strings.ToLower(strings.TrimSpace(*req.Status))
			statusSlice, ok := allowedTransport[oldStatusTask]
			if !ok {
				s.logger.Error("such status does not exist", "task_status_current", task.Status)
				return errors.New("the task status changes only in a certain order")
			}

			if !slices.Contains(statusSlice, newStatusTask) {
				s.logger.Error("can't skip status", "task_status_current", task.Status, "req_tas_status", req.Status)
				return errors.New("the task status changes only in a certain order")
			}
		}

		if req.Users != nil && task.LimitUser < len(*req.Users) {
			s.logger.Error("the number of users exceeds the allowed limit", "limit", task.LimitUser, "users_count", len(*req.Users))
			return errors.New("the number of users exceeds the allowed limit")
		}

		// Очищаем название от восклицательных знаков перед сохранением
		// Восклицательные знаки добавляются только при отображении в buildTaskResponse
		var cleanTitle *string
		if req.Title != nil {
			// Убираем все восклицательные знаки из конца названия
			cleaned := strings.TrimRight(*req.Title, "!")
			cleanTitle = &cleaned
		}

		updateReq := models.TaskUpdateReq{
			Title:       cleanTitle,
			Description: req.Description,
			Status:      req.Status,
			Users:       req.Users,
			Priority:    req.Priority,
			StartTask:   req.StartTask,
			FinishTask:  req.FinishTask,
		}

		// Управление временными метками в зависимости от изменения статуса
		if req.Status != nil {
			if newStatusTask == "in_progress" && !strings.EqualFold(oldStatusTask, "in_progress") {
				now := time.Now()
				updateReq.StartTask = &now
				updateReq.FinishTask = nil
			}
			if newStatusTask == "done" && !strings.EqualFold(oldStatusTask, "done") {
				now := time.Now()
				updateReq.FinishTask = &now
			}

			if strings.EqualFold(oldStatusTask, "done") && newStatusTask == "in_progress" {
				updateReq.FinishTask = nil
			}
			if newStatusTask == "todo" {
				updateReq.StartTask = nil
				updateReq.FinishTask = nil
			}
		}

		if err := taskrepo.UpdateTask(task.ID, updateReq); err != nil {
			s.logger.Error("failed update task from req", "err", err)
			return err
		}

		// Обновляем связь с пользователями, если они переданы
		if req.Users != nil {
			var taskModel models.Task
			if err := tx.First(&taskModel, task.ID).Error; err != nil {
				s.logger.Error("failed to get task for user update", "err", err)
				return err
			}

			oldUsersCount := len(task.Users)
			newUsersCount := len(*req.Users)

			if err := tx.Model(&taskModel).Association("Users").Replace(req.Users); err != nil {
				s.logger.Error("failed to update task users", "err", err)
				return err
			}
			s.logger.Info("task users updated", "task_id", task.ID, "old_count", oldUsersCount, "new_count", newUsersCount)

			// Логика изменения статуса при назначении/снятии пользователя
			// Если задача была "todo" и добавили пользователя - меняем на "in_progress"
			// Если задача была "in_progress" и убрали всех пользователей - меняем на "todo"
			if strings.ToLower(strings.TrimSpace(task.Status)) == "todo" && newUsersCount > 0 && oldUsersCount == 0 {
				statusInProgress := "in_progress"
				updateReq.Status = &statusInProgress
				now := time.Now()
				updateReq.StartTask = &now
				s.logger.Info("task status changed to in_progress (user assigned)", "task_id", task.ID)
			} else if strings.ToLower(strings.TrimSpace(task.Status)) == "in_progress" && newUsersCount == 0 && oldUsersCount > 0 {
				statusTodo := "todo"
				updateReq.Status = &statusTodo
				updateReq.StartTask = nil
				s.logger.Info("task status changed to todo (all users unassigned)", "task_id", task.ID)
			}

			// Если статус изменился, обновляем задачу еще раз
			if updateReq.Status != nil {
				if err := taskrepo.UpdateTask(task.ID, updateReq); err != nil {
					s.logger.Error("failed to update task status after user change", "err", err)
					return err
				}
			}
		}

		// Обновляем статус проекта на основе количества завершенных задач
		// Важно: подсчет происходит ПОСЛЕ обновления задачи, чтобы учесть новый статус
		if req.Status != nil {
			var doneTasksCount int64
			var totalTasksCount int64

			// Подсчитываем количество завершенных задач ПОСЛЕ обновления
			// Теперь текущая задача уже имеет новый статус в БД
			if err := tx.Model(&models.Task{}).
				Where("project_id = ? AND status = ?", task.ProjectID, "done").
				Count(&doneTasksCount).Error; err == nil {

				// Подсчитываем общее количество задач в проекте
				if err := tx.Model(&models.Task{}).
					Where("project_id = ?", task.ProjectID).
					Count(&totalTasksCount).Error; err == nil {

					statusDone := "done"
					statusInProgress := "in_progress"
					newProjStatus := models.ProjectUpdReq{}

					s.logger.Info("checking project status", "project_id", task.ProjectID,
						"done_tasks", doneTasksCount, "total_tasks", totalTasksCount,
						"new_task_status", newStatusTask)

					// Если все задачи завершены, проект помечаем как done
					if totalTasksCount > 0 && doneTasksCount == totalTasksCount {
						newProjStatus.Status = &statusDone
						s.logger.Info("all tasks done, setting project to done", "project_id", task.ProjectID)
					} else if doneTasksCount == 0 {
						// Если нет завершенных задач, проект в процессе
						newProjStatus.Status = &statusInProgress
						s.logger.Info("no tasks done, setting project to in_progress", "project_id", task.ProjectID)
					} else if doneTasksCount > 0 && doneTasksCount < totalTasksCount {
						// Если есть завершенные задачи, но не все - проект в процессе
						newProjStatus.Status = &statusInProgress
						s.logger.Info("some tasks done, setting project to in_progress", "project_id", task.ProjectID,
							"done", doneTasksCount, "total", totalTasksCount)
					}

					if newProjStatus.Status != nil {
						// Обновляем проект напрямую через репозиторий в транзакции
						// Валидация статуса проекта пропускается, т.к. статус определяется автоматически
						// на основе состояния задач
						if err := projectRepo.UpdateProject(task.ProjectID, newProjStatus); err != nil {
							s.logger.Error("failed to update project status", "project_id", task.ProjectID,
								"new_status", *newProjStatus.Status, "err", err)
							// Возвращаем ошибку, т.к. обновление статуса проекта важно
							return err
						} else {
							s.logger.Info("project status updated successfully", "project_id", task.ProjectID,
								"status", *newProjStatus.Status, "done_tasks", doneTasksCount, "total_tasks", totalTasksCount)
						}
					} else {
						s.logger.Info("project status unchanged", "project_id", task.ProjectID,
							"done_tasks", doneTasksCount, "total_tasks", totalTasksCount)
					}
				} else {
					s.logger.Error("failed to count total tasks", "project_id", task.ProjectID, "err", err)
					return err
				}
			} else {
				s.logger.Error("failed to count done tasks", "project_id", task.ProjectID, "err", err)
				return err
			}
		}

		s.logger.Info("update task from req successful", "op", "service.project.UpdateTask")
		return nil
	})
}

func buildTaskResponse(task *models.Task) *models.TaskResponse {
	if task == nil {
		return nil
	}

	resp := &models.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		ProjectID:   task.ProjectID,
		Users:       task.Users,
		LimitUser:   task.LimitUser,
		StartTask:   task.StartTask,
	}

	if task.FinishTask != nil {
		resp.FinishTask = task.FinishTask
	}

	switch task.Priority {
	case 1:
		resp.Title += "!"
		resp.Priority = "Важно!"
	case 2:
		resp.Title += "!!!"
		resp.Priority = "Очень важно!!!"
	default:
		resp.Priority = "Обычная"
	}

	return resp
}

func (s *taskService) AssignTaskToUser(taskID uint, userID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		taskrepo := s.repo.WithDB(tx)

		task, err := taskrepo.GetTaskByID(taskID)
		if err != nil {
			s.logger.Error("failed to get task", "task_id", taskID, "err", err)
			return err
		}

		// Проверяем, не назначен ли уже пользователь
		for _, user := range task.Users {
			if user.ID == userID {
				s.logger.Info("user already assigned to task", "task_id", taskID, "user_id", userID)
				return nil // Уже назначен
			}
		}

		// Добавляем пользователя
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			s.logger.Error("failed to get user", "user_id", userID, "err", err)
			return err
		}

		if err := tx.Model(task).Association("Users").Append(&user); err != nil {
			s.logger.Error("failed to assign user to task", "err", err)
			return err
		}

		// Если задача была "todo" и теперь есть пользователи - меняем на "in_progress"
		if strings.ToLower(strings.TrimSpace(task.Status)) == "todo" {
			statusInProgress := "in_progress"
			updateReq := models.TaskUpdateReq{
				Status: &statusInProgress,
			}
			now := time.Now()
			updateReq.StartTask = &now

			if err := taskrepo.UpdateTask(taskID, updateReq); err != nil {
				s.logger.Error("failed to update task status", "err", err)
				return err
			}
			s.logger.Info("task status changed to in_progress (user assigned)", "task_id", taskID, "user_id", userID)
		}

		return nil
	})
}

func (s *taskService) UnassignTaskFromUser(taskID uint, userID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		taskrepo := s.repo.WithDB(tx)

		task, err := taskrepo.GetTaskByID(taskID)
		if err != nil {
			s.logger.Error("failed to get task", "task_id", taskID, "err", err)
			return err
		}

		// Удаляем пользователя
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			s.logger.Error("failed to get user", "user_id", userID, "err", err)
			return err
		}

		if err := tx.Model(task).Association("Users").Delete(&user); err != nil {
			s.logger.Error("failed to unassign user from task", "err", err)
			return err
		}

		// Перезагружаем задачу чтобы получить актуальный список пользователей
		task, err = taskrepo.GetTaskByID(taskID)
		if err != nil {
			s.logger.Error("failed to reload task", "err", err)
			return err
		}

		// Если задача была "in_progress" и больше нет пользователей - меняем на "todo"
		if strings.ToLower(strings.TrimSpace(task.Status)) == "in_progress" && len(task.Users) == 0 {
			statusTodo := "todo"
			updateReq := models.TaskUpdateReq{
				Status: &statusTodo,
			}
			updateReq.StartTask = nil

			if err := taskrepo.UpdateTask(taskID, updateReq); err != nil {
				s.logger.Error("failed to update task status", "err", err)
				return err
			}
			s.logger.Info("task status changed to todo (all users unassigned)", "task_id", taskID, "user_id", userID)
		}

		return nil
	})
}
