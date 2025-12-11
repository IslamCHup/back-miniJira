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
	return s.db.Transaction(func(tx *gorm.DB) error {
		taskrepo := s.repo.WithDB(tx)
		projectRepo := s.projectRepo.WithDB(tx)

		task, err := taskrepo.GetTaskByID(id)
		if err != nil {
			s.logger.Error("failed to get the task by id",
				"op", "service.project.Update", "id", id, "err", err)
			return err
		}

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

		updateReq := models.TaskUpdateReq{
			Title:       req.Title,
			Description: req.Description,
			Status:      req.Status,
			Users:       req.Users,
			Priority:    req.Priority,
			StartTask:   req.StartTask,
			FinishTask:  req.FinishTask,
		}

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

		countWithStatus, err := taskrepo.CountTasksByStatusByProjectID(task.ProjectID, task.ID, "done")
		if err != nil {
			s.logger.Error("failed to count open tasks", "op", "service.task.UpdateTask", "project_id", task.ProjectID, "err", err)
			return err
		}

		statusDone := "done"
		statusInProgress := "in_progress"
		newProjStatus := models.ProjectUpdReq{}
		if countWithStatus == 0 {
			newProjStatus.Status = &statusInProgress
		} else {
			newProjStatus.Status = &statusDone
		}

		if err := projectRepo.UpdateProject(task.ProjectID, newProjStatus); err != nil {
			s.logger.Info("project closed (all tasks done)", "project_id", task.ProjectID)
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
