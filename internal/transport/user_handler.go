package transport

import (
	"net/http"
	"strconv"
	"log/slog"

	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
	logger  *slog.Logger
}

func NewUserHandler(service service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error("invalid user id", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.service.GetUserByID(uint(userID))
	if err != nil {
		h.logger.Error("GetUserByID failed", "id", userID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("GetUserByID success", "id", userID)
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	h.logger.Info("CreateUser handler called")

	var req models.UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("CreateUser invalid body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.CreateUser(req)
	if err != nil {
		h.logger.Error("CreateUser failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("CreateUser success", "id", user.ID)
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	h.logger.Info("UpdateUser handler called")

	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error("UpdateUser invalid user id", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req models.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("UpdateUser invalid body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)

	if err := h.service.UpdateUser(uint(userID), req, currentUser); err != nil {
		h.logger.Error("UpdateUser failed", "id", userID, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("UpdateUser success", "id", userID)
	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	h.logger.Info("DeleteUser handler called")

	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error("DeleteUser invalid user id", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)

	if err := h.service.DeleteUser(uint(userID), currentUser); err != nil {
		h.logger.Error("DeleteUser failed", "id", userID, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("DeleteUser success", "id", userID)
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
