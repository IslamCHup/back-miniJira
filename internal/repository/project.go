package repository

import (
	"back-minijira-petproject1/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	CreateProject(req *models.ProjectCreateReq) error
	GetProjectByID(id uint) (models.ProjectCreateResponse, error)
	ListProjects() ([]models.ProjectCreateResponse, error)
	UpdateProject(id uint, req models.ProjectUpdReq) error
	DeleteProject(id uint) error
}

type projectRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewProjectRepository(db *gorm.DB, logger *slog.Logger) ProjectRepository {
	return &projectRepository{db: db, logger: logger}
}

func (r *projectRepository) CreateProject(req *models.ProjectCreateReq) error {
	project := models.Project{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}

	if req.TimeEnd != nil {
		project.TimeEnd = *req.TimeEnd
	}

	res := r.db.Create(project)
	if res.Error != nil {
		r.logger.Error("create project failed", "err", res.Error)
		return res.Error
	}
	r.logger.Info("CreateProject success", "rows", res.RowsAffected)
	return nil
}

func (r *projectRepository) GetProjectByID(id uint) (models.ProjectCreateResponse, error) {
	var project models.ProjectCreateResponse
	if err := r.db.Model(&models.Project{}).First(&project, id).Error; err != nil {
		r.logger.Error("GetProjectByID failed", "id", id, "err", err)
		return models.ProjectCreateResponse{}, err
	}

	resp := models.ProjectCreateResponse{
		Title:       project.Title,
		Description: project.Description,
		Status:      project.Status,
		TimeEnd:     project.TimeEnd,
	}

	r.logger.Info("GetProjectByID success", "id", id)
	return resp, nil
}

func (r *projectRepository) ListProjects() ([]models.ProjectCreateResponse, error) {
	var projects []models.Project
	res := r.db.Find(&projects)

	if res.Error != nil {
		r.logger.Error("ListProjects failed", "err", res.Error)
		return nil, res.Error
	}
	projectResp := []models.ProjectCreateResponse{}
	for _, v := range projects {
		dto := models.ProjectCreateResponse{
			Title:       v.Title,
			Description: v.Description,
			Status:      v.Status,
			TimeEnd:     v.TimeEnd,
		}

		projectResp = append(projectResp, dto)
	}
	r.logger.Info("ListProjects success", "op", "repo.project", "count", len(projects))
	return projectResp, nil
}

func (r *projectRepository) UpdateProject(id uint, req models.ProjectUpdReq) error {
	res := r.db.Model(&models.Project{}).Where("id = ?", id).Updates(req)
	if res.Error != nil {
		r.logger.Error("UpdateProject failed", "id", id, "err", res.Error)
		return res.Error
	}
	r.logger.Info("UpdateProject success", "id", id, "rows", res.RowsAffected)
	return nil
}

func (r *projectRepository) DeleteProject(id uint) error {
	res := r.db.Delete(&models.Project{}, id)
	if res.Error != nil {
		r.logger.Error("DeleteProject failed", "id", id, "err", res.Error)
		return res.Error
	}
	r.logger.Info("DeleteProject success", "id", id, "rows", res.RowsAffected)
	return nil
}
