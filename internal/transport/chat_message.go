package handler

import (
	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/service"
	"context"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService service.ChatService
}

func NewChatHandler(chatService service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
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
	chatID, _ := strconv.Atoi(chatIDStr)

	var input struct {
		UserID uint   `json:"user_id"`
		Text   string `json:"text"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := h.chatService.AddMessage(
		context.Background(),
		models.ChatMessageCreateReq{
			UserID:       input.UserID,
			Text:         input.Text,
			ChatableID:   uint(chatID),
			ChatableType: chatType,
		},
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, msg)
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	chatType := c.Param("type")
	chatIDStr := c.Param("id")
	chatID, _ := strconv.Atoi(chatIDStr)

	messages, err := h.chatService.GetMessages(
		context.Background(),
		chatType,
		uint(chatID),
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, messages)
}
