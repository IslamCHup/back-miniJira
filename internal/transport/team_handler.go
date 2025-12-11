package transport

import (
	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	service service.TeamService
	logger  *slog.Logger
}

func NewTeamHandler(service service.TeamService, logger *slog.Logger) *TeamHandler {
	return &TeamHandler{service: service, logger: logger}
}

func (h *TeamHandler) Create(c *gin.Context) {
	var req models.TeamCreateReq

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("CreateTeam invalid JSON", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)

	resp, err := h.service.Create(req, currentUser)
	if err != nil {
		h.logger.Error("CreateTeam failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("CreateTeam success", "team_id", resp.ID)
	c.JSON(http.StatusCreated, resp)
}

func (h *TeamHandler) Update(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("id"))

	var req models.TeamUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("UpdateTeam invalid JSON", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)

	if err := h.service.Update(uint(teamID), req, currentUser); err != nil {
		h.logger.Error("UpdateTeam failed", "team_id", teamID, "err", err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("UpdateTeam success", "team_id", teamID)
	c.JSON(http.StatusOK, gin.H{"message": "team updated"})
}

func (h *TeamHandler) Delete(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("id"))
	currentUser := c.MustGet("currentUser").(models.User)

	if err := h.service.Delete(uint(teamID), currentUser); err != nil {
		h.logger.Error("DeleteTeam failed", "team_id", teamID, "err", err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("DeleteTeam success", "team_id", teamID)
	c.JSON(http.StatusOK, gin.H{"message": "team deleted"})
}

func (h *TeamHandler) GetByID(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.GetByID(uint(teamID))
	if err != nil {
		h.logger.Error("GetTeamByID failed", "team_id", teamID, "err", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("GetTeamByID success", "team_id", teamID)
	c.JSON(http.StatusOK, resp)
}
