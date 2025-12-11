package transport

import (
	"back-minijira-petproject1/internal/middleware"
	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/service"
	"log/slog"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService service.ChatService
	logger      *slog.Logger
}

func NewChatHandler(chatService service.ChatService, logger *slog.Logger) *ChatHandler {
	return &ChatHandler{chatService: chatService, logger: logger}
}

func (h *ChatHandler) SetupChatRoutes(r *gin.Engine, authService service.AuthService) {
	authChat := r.Group("/chat/:type/:id") // type = "projects" | "tasks"
	authChat.Use(middleware.AuthMiddleware(authService))
	{
		authChat.POST("/", h.AddMessage)
		authChat.GET("/", h.GetMessages)
	}

}

func (h *ChatHandler) AddMessage(c *gin.Context) {
	chatType := c.Param("type") // "projects" / "tasks"
	chatIDStr := c.Param("id")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID format"})
		return
	}

	// Валидируем тип чата
	if chatType != "projects" && chatType != "tasks" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat type, must be 'projects' or 'tasks'"})
		return
	}

	// Создаем структуру для валидации только user_id и text
	var input struct {
		UserID uint   `json:"user_id" binding:"required"`
		Text   string `json:"text" binding:"required,min=1,max=5000"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("invalid chat message body", "op", "ChatHandler.AddMessage", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	// Проверяем доступ для задач
	if chatType == "tasks" {
		hasAccess, err := h.chatService.CanUserAccessTask(c.Request.Context(), uint(chatID), input.UserID)
		if err != nil {
			h.logger.Error("ошибка при проверке прав доступа", "op", "ChatHandler.AddMessage", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access"})
			return
		}
		if !hasAccess {
			h.logger.Warn("попытка добавить комментарий без доступа", "user_id", input.UserID, "task_id", chatID)
			c.JSON(http.StatusForbidden, gin.H{"error": "you don't have access to this task"})
			return
		}
	}

	// Создаем полную структуру для сервиса
	req := models.ChatMessageCreateReq{
		UserID:       input.UserID,
		Text:         input.Text,
		ChatableID:   uint(chatID),
		ChatableType: chatType,
	}

	msg, err := h.chatService.AddMessage(c.Request.Context(), req)

	if err != nil {
		h.logger.Error("failed to add message", "op", "ChatHandler.AddMessage", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("message added", "op", "ChatHandler.AddMessage", "chatable_type", chatType, "chatable_id", chatID)
	c.JSON(http.StatusCreated, msg)
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	chatType := c.Param("type")
	chatIDStr := c.Param("id")

	h.logger.Info("GetMessages called", "op", "ChatHandler.GetMessages", "type", chatType, "id", chatIDStr)

	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		h.logger.Warn("invalid chat ID format", "op", "ChatHandler.GetMessages", "id", chatIDStr, "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID format"})
		return
	}

	// Валидируем тип чата
	if chatType != "projects" && chatType != "tasks" {
		h.logger.Warn("invalid chat type", "op", "ChatHandler.GetMessages", "type", chatType)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat type, must be 'projects' or 'tasks'"})
		return
	}

	messages, err := h.chatService.GetMessages(
		c.Request.Context(),
		chatType,
		uint(chatID),
	)

	if err != nil {
		h.logger.Error("failed to get messages", "op", "ChatHandler.GetMessages", "error", err, "type", chatType, "id", chatID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("messages retrieved", "op", "ChatHandler.GetMessages", "count", len(messages), "type", chatType, "id", chatID)
	c.JSON(http.StatusOK, messages)
}
