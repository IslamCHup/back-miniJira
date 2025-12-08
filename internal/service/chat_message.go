package service

import (
	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/repository"
	"context"
)

type ChatService interface {
	AddMessage(ctx context.Context, input models.ChatMessageCreateReq) (*models.ChatMessage, error)
	GetMessages(ctx context.Context, chatableType string, chatableID uint) ([]models.ChatMessage, error)
}

type chatService struct {
	repo repository.ChatRepository
}

func NewChatService(repo repository.ChatRepository) ChatService {
	return &chatService{repo: repo}
}

func (s *chatService) AddMessage(ctx context.Context, input models.ChatMessageCreateReq) (*models.ChatMessage, error) {
	msg := &models.ChatMessage{
		UserID:       input.UserID,
		Text:         input.Text,
		ChatableID:   input.ChatableID,
		ChatableType: input.ChatableType,
	}

	err := s.repo.Create(ctx, msg)
	return msg, err
}

func (s *chatService) GetMessages(ctx context.Context, chatableType string, chatableID uint) ([]models.ChatMessage, error) {
	return s.repo.GetByChat(ctx, chatableType, chatableID)
}
