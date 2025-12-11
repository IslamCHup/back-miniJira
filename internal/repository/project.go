package repository

import (
	"back-minijira-petproject1/internal/models"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	CreateProject(req *models.ProjectCreateReq) error
	GetProjectByID(id uint) (models.ProjectCreateResponse, error)
	ListProjects(filter *models.ProjectFilter) ([]models.ProjectCreateResponse, error)
	UpdateProject(id uint, req models.ProjectUpdReq) error
	DeleteProject(id uint) error
	WithDB(db *gorm.DB) ProjectRepository
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
		project.TimeEnd = req.TimeEnd
	}

	res := r.db.Create(&project)
	if res.Error != nil {
		r.logger.Error("create project failed", "err", res.Error)
		return res.Error
	}
	r.logger.Info("CreateProject success", "rows", res.RowsAffected)
	return nil
}

func (r *projectRepository) GetProjectByID(id uint) (models.ProjectCreateResponse, error) {
	var project models.Project
	if err := r.db.First(&project, id).Error; err != nil {
		r.logger.Error("GetProjectByID failed", "id", id, "err", err)
		return models.ProjectCreateResponse{}, err
	}

	resp := models.ProjectCreateResponse{
		ID:          project.ID,
		Title:       project.Title,
		Description: project.Description,
		Status:      project.Status,
	}

	if project.TimeEnd != nil {
		resp.TimeEnd = *project.TimeEnd
	} else {
		resp.TimeEnd = time.Time{} // Пустое время, если nil
	}

	r.logger.Info("GetProjectByID success", "id", id)
	return resp, nil
}

func (r *projectRepository) ListProjects(filter *models.ProjectFilter) ([]models.ProjectCreateResponse, error) {
	var projects []models.Project
	query := r.db.Model(models.Project{})

	if filter.Title != nil && *filter.Title != "" {
		query = query.Where("title = ?", *filter.Title)
	}
	if filter.Description != nil && *filter.Description != "" {
		query = query.Where("description LIKE ?", "%"+*filter.Description+"%")
	}
	if filter.Status != nil && *filter.Status != "" {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	res := query.Find(&projects)

	if res.Error != nil {
		r.logger.Error("ListProjects failed", "err", res.Error)
		return nil, res.Error
	}
	projectResp := []models.ProjectCreateResponse{}
	for _, v := range projects {
		dto := models.ProjectCreateResponse{
			ID:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			Status:      v.Status,
		}

		if v.TimeEnd != nil {
			dto.TimeEnd = *v.TimeEnd
		} else {
			dto.TimeEnd = time.Time{} // Пустое время, если nil
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

func (r *projectRepository) WithDB(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db, logger: r.logger}
}
