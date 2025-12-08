package transport

import (
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

func (h *ChatHandler) SetupChatRoutes(r *gin.Engine) {
	chat := r.Group("/chat/:type/:id") // type = "projects" | "tasks"
	{
		chat.POST("/", h.AddMessage)
		chat.GET("/", h.GetMessages)
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

	var input models.ChatMessageCreateReq
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка прав доступа: есть ли пользователь в задаче
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

	input.ChatableID = uint(chatID)
	input.ChatableType = chatType

	msg, err := h.chatService.AddMessage(c.Request.Context(), input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, msg)
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	chatType := c.Param("type")
	chatIDStr := c.Param("id")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID format"})
		return
	}

	messages, err := h.chatService.GetMessages(
		c.Request.Context(),
		chatType,
		uint(chatID),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}
