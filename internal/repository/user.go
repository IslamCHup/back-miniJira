package repository

import (
	"back-minijira-petproject1/internal/models"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(req *models.User) error
	GetUserByID(id uint) (models.User, []uint, error)
	UpdateUser(req *models.User) error
	DeleteUser(id uint) error
	AssignTasksToUser(user *models.User, taskIDs []uint) error
	GetUserByEmail(email string) (models.User, error)
	GetUserVerifyToken(token string) (models.User, error)
	UpdateUserVerification(id uint, isVerified bool, token string) error
	CountUsers() (int64, error)
	ListUsers() ([]models.User, error)
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
	r.logger.Info("user created", "id", req.ID)
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

func (r *userRepository) AssignTasksToUser(user *models.User, taskIDs []uint) error {
	if len(taskIDs) == 0 {
		return nil
	}

	r.logger.Info("assignTasksToUser called", "user_id", user.ID, "task_ids", taskIDs)

	var tasks []models.Task
	if err := r.db.Where("id IN ?", taskIDs).Find(&tasks).Error; err != nil {
		r.logger.Error("assignTasksToUser: invalid task IDs", "error", err)
		return errors.New("некорректное айди")
	}

	if err := r.db.Model(user).Association("Tasks").Replace(tasks); err != nil {
		r.logger.Error("assignTasksToUser: failed replacing tasks", "error", err)
		return err
	}

	return nil
}

func (r *userRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		r.logger.Error("Поиск по емайлу не удался", "email", email, "error", err)
		return models.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetUserVerifyToken(token string) (models.User, error) {
	var user models.User

	if err := r.db.Where("verify_token = ?", token).First(&user).Error; err != nil {
		r.logger.Error("")
		return models.User{}, nil
	}
	return user, nil
}

func (r *userRepository) UpdateUserVerification(id uint, isVerified bool, token string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"is_verified":  isVerified,
			"verify_token": token,
		}).Error
}

func (r *userRepository) CountUsers() (int64, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *userRepository) ListUsers() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		r.logger.Error("ListUsers failed", "err", err)
		return nil, err
	}
	r.logger.Info("ListUsers success", "count", len(users))
	return users, nil
}
