package service

import (
	"errors"
	"log/slog"

	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/repository"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(req models.UserCreateReq) (models.UserResponse, error)
	GetUserByID(id uint) (models.UserResponse, error)
	UpdateUser(id uint, req models.UserUpdateReq, currentUser models.User) error
	checkUserPermission(currentUser, targetUser models.User) error
	DeleteUser(id uint, currentUser models.User) error
	ListUsers() ([]models.UserResponse, error)
}

type userService struct {
	repo   repository.UserRepository
	db     *gorm.DB
	logger *slog.Logger
}

func NewUserService(repo repository.UserRepository, db *gorm.DB, logger *slog.Logger) UserService {
	return &userService{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

func (s *userService) CreateUser(req models.UserCreateReq) (models.UserResponse, error) {
	s.logger.Info("CreateUser called", "full_name", req.FullName)

	user := models.User{
		FullName: req.FullName,
	}

	if err := s.repo.CreateUser(&user); err != nil {
		return models.UserResponse{}, err
	}

	if err := s.repo.AssignTasksToUser(&user, req.TaskIDs); err != nil {
		return models.UserResponse{}, err
	}

	_, taskIDs, err := s.repo.GetUserByID(user.ID)
	if err != nil {
		return models.UserResponse{}, err
	}

	s.logger.Info("CreateUser success", "id", user.ID)

	return models.UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		IsAdmin:  user.IsAdmin,
		TaskIDs:  taskIDs,
	}, nil
}

func (s *userService) GetUserByID(id uint) (models.UserResponse, error) {
	s.logger.Info("GetUserByID called", "id", id)

	user, taskIDs, err := s.repo.GetUserByID(id)
	if err != nil {
		s.logger.Error("GetUserByID failed", "id", id, "error", err)
		return models.UserResponse{}, err
	}

	s.logger.Info("GetUserByID success", "id", id)

	return models.UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		IsAdmin:  user.IsAdmin,
		TaskIDs:  taskIDs,
	}, nil
}

func (s *userService) UpdateUser(id uint, req models.UserUpdateReq, currentUser models.User) error {
	user, _, err := s.repo.GetUserByID(id)
	if err != nil {
		return err
	}

	if err := s.checkUserPermission(currentUser, user); err != nil {
		return err
	}

	if req.FullName != nil {
		if *req.FullName == "" {
			return errors.New("fullname cannot be empty")
		}
		user.FullName = *req.FullName
	}

	if req.TaskIDs != nil {
		if err := s.repo.AssignTasksToUser(&user, req.TaskIDs); err != nil {
			return err
		}
	}

	return s.repo.UpdateUser(&user)
}

func (s *userService) checkUserPermission(currentUser, targetUser models.User) error {
	if currentUser.ID != targetUser.ID && !currentUser.IsAdmin {
		s.logger.Warn("permission denied",
			"current_user_id", currentUser.ID,
			"target_user_id", targetUser.ID,
		)
		return errors.New("you cannot update another user")
	}

	return nil
}

func (s *userService) DeleteUser(id uint, currentUser models.User) error {
	s.logger.Info("DeleteUser called", "target_id", id, "by_user", currentUser.ID)

	user, _, err := s.repo.GetUserByID(id)
	if err != nil {
		s.logger.Error("DeleteUser: user not found", "id", id, "error", err)
		return err
	}

	if err := s.checkUserPermission(currentUser, user); err != nil {
		return err
	}

	if err := s.repo.DeleteUser(id); err != nil {
		s.logger.Error("DeleteUser failed", "id", id, "error", err)
		return err
	}

	s.logger.Info("DeleteUser success", "id", id)
	return nil
}

func (s *userService) ListUsers() ([]models.UserResponse, error) {
	users, err := s.repo.ListUsers()
	if err != nil {
		s.logger.Error("ListUsers failed", "err", err)
		return nil, err
	}

	var responses []models.UserResponse
	for _, user := range users {
		_, taskIDs, err := s.repo.GetUserByID(user.ID)
		if err != nil {
			s.logger.Warn("ListUsers: failed to get task IDs", "user_id", user.ID, "err", err)
			taskIDs = []uint{}
		}

		responses = append(responses, models.UserResponse{
			ID:       user.ID,
			FullName: user.FullName,
			IsAdmin:  user.IsAdmin,
			TaskIDs:  taskIDs,
		})
	}

	s.logger.Info("ListUsers success", "count", len(responses))
	return responses, nil
}
