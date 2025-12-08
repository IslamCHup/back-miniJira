package service

import (
	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/repository"
	"context"
	"errors"
	"log/slog"
	"strings"
)

var (
	ErrUserIdZero       = errors.New("user_id cannot be zero")
	ErrUserIdTaskIdZero = errors.New("user_id or task_id cannot be zero")
	ErrTextEmpty        = errors.New("text cannot be empty")
	ErrChatableIdZero   = errors.New("chatable_id cannot be zero")
	ErrInvalidChatType  = errors.New("invalid chatable_type: must be 'projects' or 'tasks'")
)

type ChatService interface {
	AddMessage(ctx context.Context, input models.ChatMessageCreateReq) (*models.ChatMessage, error)
	GetMessages(ctx context.Context, chatableType string, chatableID uint) ([]models.ChatMessage, error)
	CanUserAccessTask(ctx context.Context, taskID, userID uint) (bool, error)
}

type chatService struct {
	repo   repository.ChatRepository
	logger *slog.Logger
}

func NewChatService(repo repository.ChatRepository, logger *slog.Logger) ChatService {
	return &chatService{repo: repo, logger: logger}
}

func (s *chatService) AddMessage(ctx context.Context, input models.ChatMessageCreateReq) (*models.ChatMessage, error) {
	if input.UserID == 0 {
		s.logger.Warn("попытка добавить комментарий с нулевым UserID", "op", "chatService.AddMessage")
		return nil, ErrUserIdZero
	}

	if strings.TrimSpace(input.Text) == "" {
		s.logger.Warn("попытка добавить пустой комментарий", "op", "chatService.AddMessage", "user_id", input.UserID)
		return nil, ErrTextEmpty
	}

	if input.ChatableID == 0 {
		s.logger.Warn("попытка добавить комментарий с нулевым ChatableID", "op", "chatService.AddMessage", "user_id", input.UserID)
		return nil, ErrChatableIdZero
	}

	validTypes := map[string]bool{
		"projects": true,
		"tasks":    true,
	}
	if !validTypes[input.ChatableType] {
		s.logger.Warn("попытка добавить комментарий с неверным типом", "op", "chatService.AddMessage", "type", input.ChatableType)
		return nil, ErrInvalidChatType
	}

	msg := &models.ChatMessage{
		UserID:       input.UserID,
		Text:         strings.TrimSpace(input.Text),
		ChatableID:   input.ChatableID,
		ChatableType: input.ChatableType,
	}

	err := s.repo.Create(ctx, msg)
	if err != nil {
		s.logger.Error("ошибка при создании комментария", "op", "chatService.AddMessage", "error", err)
		return nil, err
	}

	s.logger.Info("комментарий успешно создан", "op", "chatService.AddMessage", "message_id", msg.ID, "user_id", input.UserID)
	return msg, nil
}

func (s *chatService) GetMessages(ctx context.Context, chatableType string, chatableID uint) ([]models.ChatMessage, error) {
	if chatableID == 0 {
		s.logger.Warn("попытка получить комментарии с нулевым ID", "op", "chatService.GetMessages", "type", chatableType)
		return nil, ErrChatableIdZero
	}

	messages, err := s.repo.GetByChat(ctx, chatableType, chatableID)
	if err != nil {
		s.logger.Error("ошибка при получении комментариев", "op", "chatService.GetMessages", "error", err, "type", chatableType, "id", chatableID)
		return nil, err
	}

	s.logger.Info("комментарии успешно получены", "op", "chatService.GetMessages", "count", len(messages), "type", chatableType, "id", chatableID)
	return messages, nil
}

func (s *chatService) CanUserAccessTask(ctx context.Context, taskID, userID uint) (bool, error) {
	if taskID == 0 || userID == 0 {
		return false, ErrUserIdTaskIdZero
	}

	hasAccess, err := s.repo.IsUserInTask(ctx, taskID, userID)
	if err != nil {
		s.logger.Error("ошибка при проверке доступа пользователя", "op", "chatService.CanUserAccessTask", "error", err, "task_id", taskID, "user_id", userID)
		return false, err
	}

	if !hasAccess {
		s.logger.Warn("пользователь не имеет доступа к задаче", "op", "chatService.CanUserAccessTask", "task_id", taskID, "user_id", userID)
	}

	return hasAccess, nil
}
