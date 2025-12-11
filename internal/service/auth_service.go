package service

import (
	"back-minijira-petproject1/internal/auth"
	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/repository"

	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req models.RegisterRequest) error
	Login(req models.LoginRequest) (string, error)
	GetUserByID(id uint) (models.User, error)
	VerifyEmail(token string) error
}

type authService struct {
	repo         repository.UserRepository
	logger       *slog.Logger
	emailService *EmailService
}

func NewAuthService(repo repository.UserRepository, logger *slog.Logger) AuthService {
	return &authService{repo: repo, logger: logger, emailService: NewEmailService()}
}

func (s *authService) Register(req models.RegisterRequest) error {
	if exisiting, _ := s.repo.GetUserByEmail(req.Email); exisiting.ID != 0 {
		return errors.New("пользователь с таким емайлом уже есть")

	}

	  count, err := s.repo.CountUsers()
    if err != nil {
        return err
    }


	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	verifyToken := uuid.New().String()

	user := models.User{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: string(hash),
		IsVerified:   false,
		VerifyToken:  verifyToken,
		IsAdmin:      count == 0,
	}

	if err := s.repo.CreateUser(&user); err != nil {
		return err
	}

	verifyLink := fmt.Sprintf("http://localhost:8080/auth/verify?token=%s", verifyToken)
	if err := s.emailService.SendVerificationEmail(user.Email, user.FullName, verifyLink); err != nil {
		s.logger.Error("email send failed", "error", err)
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

func (s *authService) Login(req models.LoginRequest) (string, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !user.IsVerified {
		return "", errors.New("email is not verified")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	return auth.GenerateToken(user.ID, user.IsAdmin)
}

func (s *authService) VerifyEmail(token string) error {
	user, err := s.repo.GetUserVerifyToken(token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	user.IsVerified = true
	user.VerifyToken = ""

	return s.repo.UpdateUserVerification(user.ID, user.IsVerified, user.VerifyToken)
}

func (s *authService) GetUserByID(id uint) (models.User, error) {
	user, _, err := s.repo.GetUserByID(id)
	return user, err
}
