package transport

import (
	"back-minijira-petproject1/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	service service.ReportService
}

func NewReportHandler(s service.ReportService) *ReportHandler {
	return &ReportHandler{service: s}
}

func (h *ReportHandler) RegisterRoutes(r *gin.Engine) {
	reports := r.Group("/projects/:id/reports")
	{
		reports.GET("/top-workers")
		reports.GET("avg-time")
		reports.GET("completion-percent")
		reports.GET("/user-tracker/:userId")
	}
}

func (h *ReportHandler) GetTopWorkers(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.service.TopWorkers(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, data)
}

func (h *ReportHandler) GetAverageTime(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.service.AverageTime(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, data)
}

func (h *ReportHandler) GetCompletionPercent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.service.CompletionPercent(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, data)
}

func (h *ReportHandler) GetUserTracker(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := strconv.Atoi(c.Param("userId"))

	data, err := h.service.UserTracker(uint(projectID), uint(userID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, data)
}
