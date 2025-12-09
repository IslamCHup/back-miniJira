package repository

import (
	"back-minijira-petproject1/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(req *models.User) error
	GetUserByID(id uint) (models.User, []uint, error)
	UpdateUser(req *models.User) error
	DeleteUser(id uint) error
}

type userRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewUserRepository(db *gorm.DB, logger *slog.Logger) UserRepository {
	return &userRepository{db: db, logger: logger}
}

func (r *userRepository) CreateUser(req *models.User) error {
	if err := r.db.Create(req).Error; err != nil {
		r.logger.Error("failed to create user", "error", err)
		return err
	}
	r.logger.Info("user created","id",req.ID)
	return nil
}

func (r *userRepository) GetUserByID(id uint) (models.User, []uint, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		r.logger.Error("GetByID failed", "id", id, "error", err)
		return user, nil, err
	}

	var taskIDs []uint
	if err := r.db.Table("user_tasks").
		Where("user_id = ?", id).
		Pluck("task_id", &taskIDs).Error; err != nil {
		r.logger.Error("failed to load task IDs for user", "id", id, "error", err)
		return user, nil, err
	}
	r.logger.Info("GetByID success", "id", id)
	return user, taskIDs, nil
}

func (r *userRepository) UpdateUser(req *models.User) error {
	res := r.db.Save(req)
	if res.Error != nil {
		r.logger.Error("UpdateUser failed", "id", req.ID)
		return res.Error
	}
	r.logger.Info("UpdateUser succees", "id", req.ID, "rows", res.RowsAffected)
	return nil
}

func (r *userRepository) DeleteUser(id uint) error {
	res := r.db.Delete(&models.User{}, id)
	if res.Error != nil {
		r.logger.Error("DeleteUser failed", "id", id, "err", res.Error)
		return res.Error
	}
	r.logger.Info("DeleteUser success", "id", id, "rows", res.RowsAffected)
	return nil
}
