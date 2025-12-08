package repository

import (
	"back-minijira-petproject1/internal/models"
	"context"

	"gorm.io/gorm"
)

type ChatRepository interface {
	Create(ctx context.Context, msg *models.ChatMessage) error
	GetByChat(ctx context.Context, chatableType string, chatableID uint) ([]models.ChatMessage, error)
	IsUserInTask(ctx context.Context, taskID, userID uint) (bool, error)
}

type chatRepositoryGorm struct {
	db *gorm.DB
}

func NewChatRepositoryGorm(db *gorm.DB) ChatRepository {
	return &chatRepositoryGorm{db: db}
}

func (r *chatRepositoryGorm) Create(ctx context.Context, msg *models.ChatMessage) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *chatRepositoryGorm) GetByChat(ctx context.Context, chatableType string, chatableID uint) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.WithContext(ctx).
		Where("chatable_type = ? AND chatable_id = ?", chatableType, chatableID).
		Order("created_at ASC").
		Find(&messages).Error

	return messages, err
}

func (r *chatRepositoryGorm) IsUserInTask(ctx context.Context, taskID, userID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("task_users").
		Where("task_id = ? AND user_id = ?", taskID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
