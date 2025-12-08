package repository

import (
	"back-minijira-petproject1/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(req models.TaskCreateReq) error
	UpdateTask(id uint, req models.TaskUpdateReq) error
	DeleteTask(id uint) error
	ListTasks() ([]*models.Task, error)
	GetTaskByID(id uint) (*models.Task, error)
}

type taskRepository struct {
	db     *gorm.DB
	logger slog.Logger
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
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

func (r *taskRepository) ListTasks() ([]*models.Task, error) {
	var tasks []*models.Task
	if err := r.db.Find(&tasks).Error; err != nil {
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
