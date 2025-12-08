package repository

import (
	"back-minijira-petproject1/models"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	CreateProject(req *models.ProjectCreateReq) error
	GetProjectByID(id uint) (models.Project, error)
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
	res := r.db.Create(&req)
	if res.Error != nil {
		r.logger.Error("create project failed", "err", res.Error)
		return res.Error
	}
	r.logger.Info("CreateProject success", "rows", res.RowsAffected)
	return nil
}

func (r *projectRepository) GetProjectByID(id uint) (models.Project, error) {
	var project models.Project
	if err := r.db.First(&project, id).Error; err != nil {
		r.logger.Error("GetProjectByID failed", "id", id, "err", err)
		return models.Project{}, err
	}
	r.logger.Info("GetProjectByID success", "id", id)
	return project, nil
}

func (r *projectRepository) ListProjects() ([]models.ProjectCreateResponse, error) {
	var projects []models.ProjectCreateResponse
	res := r.db.Model(&models.Project{}).Find(&projects)

	if res.Error != nil {
		r.logger.Error("ListProjects failed", "err", res.Error)
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		r.logger.Error("ListProjects empty", "op", "repo.project", "rows", res.RowsAffected)
		return nil, errors.New("ListProjects empty")
	}

	r.logger.Info("ListProjects success", "op", "repo.project", "count", len(projects))
	return projects, nil
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
