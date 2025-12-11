package transport

import (
	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	service service.ProjectService
	logger  *slog.Logger
}

func NewProjectHandler(service service.ProjectService, logger *slog.Logger) ProjectHandler {
	return ProjectHandler{service: service, logger: logger}
}

// ...existing code...
func (h *ProjectHandler) RegisterRoutes(r *gin.Engine) {
    projects := r.Group("/projects")
    {
        projects.GET("/", h.ListProjects)
        projects.GET("/:id", h.GetByID)
        projects.POST("/", h.Create)
        projects.PATCH("/:id", h.UpdateProject)
        projects.DELETE("/:id", h.Delete)
    }
}
// ...existing code...

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	title := c.Query("title")
	description := c.Query("description")
	status := c.Query("status")
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	filter := models.ProjectFilter{
		Title:       &title,
		Description: &description,
		Status:      &status,
		Limit:       limit,
		Offset:      offset,
	}

	projects, err := h.service.ListProjects(&filter)
	if err != nil {
		h.logger.Error("failed to get projects list", "op", "handler.ListProjects", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get projects list"})
		return
	}

	h.logger.Info("projects list returned", "op", "handler.ListProjects", "count", len(projects))
	c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	project, err := h.service.GetByID(uint(id))

	if err != nil {
		h.logger.Error("failed to get project by id", "op", "handler.GetByID", "id", id, "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get project by id"})
		return
	}

	h.logger.Info("project retrieved", "op", "handler.GetByID", "id", id)
	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var projectInput models.ProjectCreateReq

	if err := c.ShouldBindJSON(&projectInput); err != nil {
		h.logger.Error("failed to parse JSON body", "op", "handler.Create", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	project, err := h.service.Create(&projectInput)

	if err != nil {
		h.logger.Error("failed to create project", "op", "handler.Create", "err", err, "project_input", projectInput)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create project"})
		return
	}

	h.logger.Info("project created", "op", "handler.Create", "project", project)
	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := h.service.Delete(uint(id)); err != nil {
		h.logger.Error("failed to delete project by id", "op", "handler.Delete", "id", id, "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to delete project by id"})
		return
	}

	h.logger.Info("project deleted", "op", "handler.Delete", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "delete successful"})
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var req models.ProjectUpdReq

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to parse JSON body", "op", "handler.Create", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.UpdateProject(uint(id), req); err != nil {
		h.logger.Error("failed to update project by id", "op", "handler.update", "err", err, "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update successful"})
}
