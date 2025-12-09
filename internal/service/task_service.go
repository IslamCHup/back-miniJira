package service

import (
	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/repository"
	"errors"
	"log/slog"
	"slices"

	"gorm.io/gorm"
)

type TaskService interface {
	GetTaskByID(id uint) (*models.TaskResponse, error)
	ListTasks(filter *models.TaskFilter) ([]*models.TaskResponse, error)
	DeleteTask(id uint) error
	CreateTask(req models.TaskCreateReq) error
	UpdateTask(id uint, req models.TaskUpdateReq) error
}

type taskService struct {
	db     *gorm.DB
	logger *slog.Logger
	repo   repository.TaskRepository
}

func NewTaskService(db *gorm.DB, logger *slog.Logger, repo repository.TaskRepository) TaskService {
	return &taskService{db: db, logger: logger, repo: repo}
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

func (s *taskService) CreateTask(req models.TaskCreateReq) error {
	if err := s.repo.CreateTask(req); err != nil {
		s.logger.Error("failed create task from req", "err", err, "req", req)
		return err
	}
	s.logger.Info("create task from req successful", "op", "service.project.CreateTask")
	return nil
}

func (s *taskService) UpdateTask(id uint, req models.TaskUpdateReq) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		taskrepo := s.repo.WithDB(tx)
		task, err := taskrepo.GetTaskByID(id)

		if err != nil {
			s.logger.Error("failed to get the task by id",
				"op", "service.project.Update", "id", id, "err", err)
			return err
		}

		allowedTransport := map[string][]string{
			"todo":        {"in_progress"},
			"in_progress": {"todo", "done"},
			"done":        {"in_progress"},
		}

		if req.Status == nil {
			s.logger.Error("such status does not exist", "task_status_current", task.Status)
			return errors.New("the task status changes only in a certain order")
		}

		statusSlice, ok := allowedTransport[task.Status]

		if !ok {
			s.logger.Error("such status does not exist", "task_status_current", task.Status)
			return errors.New("the task status changes only in a certain order")
		}

		if !slices.Contains(statusSlice, *req.Status) {
			s.logger.Error("can't skip status", "task_status_current", task.Status, "req_tas_status", req.Status)
			return errors.New("the task status changes only in a certain order")
		}

		if req.Users != nil && task.LimitUser < len(req.Users) {
			s.logger.Error("the number of users exceeds the allowed limit", "limit", task.LimitUser, "users_count", len(req.Users))
			return errors.New("the number of users exceeds the allowed limit")
		}

		if err := taskrepo.UpdateTask(id, req); err != nil {
			s.logger.Error("failed update task from req", "err", err)
			return err
		}
		s.logger.Info("update task from req successful", "op", "service.project.UpdateTask")
		return nil
	})
}

func buildTaskResponse(task *models.Task) *models.TaskResponse {
	taskResponse := models.TaskResponse{
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		ProjectID:   task.ProjectID,
		Users:       task.Users,
		LimitUser:   task.LimitUser,
		StartTask:   task.StartTask,
		FinishTask:  task.FinishTask,
	}

	switch task.Priority {
	case 1:
		taskResponse.Title += "!"
		taskResponse.Priority = "Важно!"
	case 2:
		taskResponse.Title += "!!!"
		taskResponse.Priority = "Очень важно!!!"
	}
	return &taskResponse
}
