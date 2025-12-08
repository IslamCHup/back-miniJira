package repository

import (
	"back-minijira-petproject1/models"
	"log/slog"

	"gorm.io/gorm"
)

type TaskCommentRepo interface {
	CreateTaskComment(req *models.TaskCommentCreateReq) error

	GetTaskCommentsByTaskID(taskID uint) ([]models.TaskComment, error)

	UpdateTaskComment(req *models.TaskCommentUpdateReq) error

	DeleteTaskComment(id uint) error

	GetByID(id uint) (*models.TaskComment, error)
}

type taskCommentRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewTaskCommentRepository(db *gorm.DB, logger *slog.Logger) TaskCommentRepo {
	return &taskCommentRepository{
		db:     db,
		logger: logger,
	}
}

func (r *taskCommentRepository) GetByID(id uint) (*models.TaskComment, error) {
	var comment models.TaskComment

	if err := r.db.Model(&models.TaskComment{}).Where("id = ?", id).First(&comment).Error; err != nil {
		r.logger.Error("GetByID failed", "id", id, "err", err)
		return nil, err
	}

	r.logger.Info("GetByID success", "id", id)
	return &comment, nil
}

func (r *taskCommentRepository) CreateTaskComment(req *models.TaskCommentCreateReq) error {
	if err := r.db.Create(&req).Error; err != nil {
		r.logger.Error("Create TaskComment falied")
		return err
	}
	r.logger.Info("Created successfull")
	return nil
}

func (r *taskCommentRepository) GetTaskCommentsByTaskID(taskID uint) ([]models.TaskComment, error) {
	var comments []models.TaskComment

	if err := r.db.Model(&models.TaskComment{}).Where("task_id = ?", taskID).Find(&comments); err != nil {
		r.logger.Error("GetTaskCommentsByTaskID failed", "taskId", taskID, "err", err)
		return []models.TaskComment{}, err.Error
	}

	r.logger.Info("GetTaskCommentsByTaskID success", "count", len(comments))
	return comments, nil
}

func (r *taskCommentRepository) UpdateTaskComment(req *models.TaskCommentUpdateReq) error {
	if req == nil {
		return nil
	}
	return r.db.Save(&req).Error
}

func (r *taskCommentRepository) DeleteTaskComment(id uint) error {
	return r.db.Delete(&models.TaskComment{}, id).Error
}
