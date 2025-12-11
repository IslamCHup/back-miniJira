package repository

import (
	"back-minijira-petproject1/internal/models"
	"log/slog"
	"strings"

	"gorm.io/gorm"
)

type TaskRepository interface {
	WithDB(db *gorm.DB) TaskRepository
	CreateTask(req *models.TaskCreateReq) error
	UpdateTask(id uint, req models.TaskUpdateReq) error
	DeleteTask(id uint) error
	ListTasks(filter *models.TaskFilter) ([]*models.Task, error)
	GetTaskByID(id uint) (*models.Task, error)
	CountTasksByStatusByProjectID(project_id uint, task_id uint, status string) (int64, error)
}

type taskRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewTaskRepository(db *gorm.DB, logger *slog.Logger) TaskRepository {
	return &taskRepository{db: db, logger: logger}
}

func (r *taskRepository) CreateTask(req *models.TaskCreateReq) error {
	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		ProjectID:   req.ProjectID,
		Users:       req.Users,
		Priority:    req.Priority,
		LimitUser:   req.LimitUser,
		StartTask:   req.StartTask,
		FinishTask:  req.FinishTask,
	}
	res := r.db.Create(&task)
	if res.Error != nil {
		r.logger.Error("CreateTask failed")
		return res.Error
	}
	r.logger.Info("CreateTask success", "rows", res.RowsAffected)
	return nil
}

func (r *taskRepository) UpdateTask(id uint, req models.TaskUpdateReq) error {
	res := r.db.Model(&models.Task{}).Where("id = ?", id).Omit("Users").Updates(req)
	if res.Error != nil {
		r.logger.Error("UpdateTask failed", "id", id, "err", res.Error)
		return res.Error
	}
	r.logger.Info("UpdateTask success", "id", id, "rows", res.RowsAffected)
	return nil
}

func (r *taskRepository) DeleteTask(id uint) error {
	res := r.db.Delete(&models.Task{}, id)
	if res.Error != nil {
		r.logger.Error("DeleteTask failed", "id", id, "err", res.Error)
		return res.Error
	}
	r.logger.Info("DeleteTask success", "id", id, "rows", res.RowsAffected)
	return nil
}

func (r *taskRepository) ListTasks(filter *models.TaskFilter) ([]*models.Task, error) {
	var tasks []*models.Task

	query := r.db.Model(&models.Task{})

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.UserID != nil {
		query = query.Joins("JOIN task_users ON task_users.task_id = tasks.id").
			Where("task_users.user_id = ?", *filter.UserID)
	}

	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}

	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}

	if filter.Search != nil {
		search := "%" + *filter.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", search, search)
	}

	sortField := "priority"
	sortOrder := "DESC"
	if filter != nil && filter.SortBy != nil && *filter.SortBy != "" {
		switch strings.ToLower(*filter.SortBy) {
		case "priority":
			sortField = "priority"
		case "created_at":
			sortField = "created_at"
		case "start_task":
			sortField = "start_task"
		default:
			sortField = "priority"
		}
	}
	if filter != nil && filter.SortOrder != nil && strings.ToLower(*filter.SortOrder) == "asc" {
		sortOrder = "ASC"
	}
	orderClause := sortField + " " + sortOrder + ", created_at DESC"
	query = query.Order(orderClause)

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	// Загружаем задачи с пользователями
	if err := query.Preload("Users").Find(&tasks).Error; err != nil {
		r.logger.Error("ListTask failed", "err", err)
		return nil, err
	}
	r.logger.Info("ListTask success", "count", len(tasks))
	return tasks, nil
}

func (r *taskRepository) GetTaskByID(id uint) (*models.Task, error) {
	var task models.Task
	if err := r.db.Preload("Users").Where("id = ?", id).First(&task).Error; err != nil {
		r.logger.Error("GetTaskByID failed", "id", id, "err", err)
		return nil, err
	}
	r.logger.Info("GetTaskByID success", "id", id)
	return &task, nil
}

func (r *taskRepository) WithDB(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db, logger: r.logger}
}

func (r *taskRepository) CountTasksByStatusByProjectID(project_id uint, task_id uint, status string) (int64, error) {
	var countStatus int64
	status = strings.ToLower(strings.TrimSpace(status))
	if err := r.db.Model(&models.Task{}).Where("project_id = ? AND id <> ? AND (status IS NULL OR LOWER(status) <>?)",
		project_id, task_id, status).Count(&countStatus).Error; err != nil {
		r.logger.Error("CountTasksByStatusByProjectID failed", "project_id", project_id, "task_id", task_id, "err", err)
		return -1, err
	}

	return countStatus, nil
}
