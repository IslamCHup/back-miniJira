package repository

import (
	"back-minijira-petproject1/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type TaskRepository interface {
	WithDB(db *gorm.DB) TaskRepository
	CreateTask(req models.TaskCreateReq) error
	UpdateTask(id uint, req models.TaskUpdateReq) error
	DeleteTask(id uint) error
	ListTasks(filter *models.TaskFilter) ([]*models.Task, error) 
	GetTaskByID(id uint) (*models.Task, error)
}

type taskRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewTaskRepository(db *gorm.DB, logger *slog.Logger) TaskRepository {
	return &taskRepository{db: db, logger: logger}
}

func (r *taskRepository) CreateTask(req models.TaskCreateReq) error {
	res := r.db.Create(&req)
	if res.Error != nil {
		r.logger.Error("CreateTask failed")
		return res.Error
	}
	r.logger.Info("CreateTask success", "rows", res.RowsAffected)
	return nil
}

func (r *taskRepository) UpdateTask(id uint, req models.TaskUpdateReq) error {
	res := r.db.Model(&models.Task{}).Where("id = ?", id).Updates(&req)
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

	if filter.Status != nil{
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.UserID != nil{
		query = query.Joins("JOIN task_users ON task_users.task_id = tasks.id").
		Where("task_users.user_id = ?", *filter.UserID)
	}

	if filter.Priority != nil{
		query = query.Where("priority = ?", *filter.Priority)
	}

	if filter.ProjectID != nil{
		query = query.Where("project_id = ?", *filter.ProjectID)
	}

	if filter.Search != nil{
		search := "%" + *filter.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", search, search)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	if err := query.Find(&tasks).Error; err != nil {
		r.logger.Error("ListTask failed", "err", err)
		return nil, err
	}
	r.logger.Info("ListTask success", "count", len(tasks))
	return tasks, nil
}

func (r *taskRepository) GetTaskByID(id uint) (*models.Task, error) {
	var task *models.Task
	if err := r.db.Find(&task).Error; err != nil {
		r.logger.Error("ListTask failed", "id", id, "err", err)
		return nil, err
	}
	r.logger.Info("ListTask success", "id", id)
	return task, nil
}

func (r *taskRepository) WithDB(db *gorm.DB) TaskRepository {
	return &taskRepository{db:db, logger: r.logger}
}


