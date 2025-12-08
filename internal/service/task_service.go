package service

import (
	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/repository"
	"log/slog"

	"gorm.io/gorm"
)

type TaskService interface {
}

type taskService struct {
	db      *gorm.DB
	logger  *slog.Logger
	service repository.TaskRepository
}

func NewTaskService(db *gorm.DB, logger *slog.Logger, service repository.TaskRepository) TaskService {
	return &taskService{db: db, logger: logger, service: service}
}

func (s *taskService) GetTaskByID(id uint) (*models.Task, error) {
	task, err := s.service.GetTaskByID(id)
	if err != nil {
		s.logger.Error("передан пустой запрос",
			"op", "service.task.GetTaskByID",
			"error", err)
		return nil, err
	}
	s.logger.Info("get task by id successful", "op", "service.project.GetTaskByID", "task", task)
	return task, nil
}

func (s *taskService) ListTasks() ([]*models.Task, error) {
	tasks, err := s.service.ListTasks()
	if err != nil {
		s.logger.Error("передан пустой запрос",
			"op", "service.task.ListTasks",
			"error", err)
		return nil, err
	}
	s.logger.Info("get task by id successful", "op", "service.project.GetTaskByID", "task", tasks)
	return tasks, nil
}

func (s *taskService) DeleteTask(id uint) error {
	if err := s.service.DeleteTask(id); err != nil {
		s.logger.Error("failed delete task by id", "id", id, "err", err)
		return err
	}
	s.logger.Info("delete task by id successful", "op", "service.project.DeleteTask")
	return nil
}

func (s *taskService) CreateTask(req models.TaskCreateReq) error{
	if err := s.service.CreateTask(req); err != nil {
		s.logger.Error("failed create task from req", "err", err)
		return err
	}
	s.logger.Info("create task from req successful", "op", "service.project.CreateTask")
	return nil
}
/*CreateTask(req models.TaskCreateReq) error
UpdateTask(id uint, req models.TaskUpdateReq) error
*/
