package service

import (
	"back-minijira-petproject1/models"
	"back-minijira-petproject1/repository"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type ProjectService interface {
}

type projectService struct {
	db      *gorm.DB
	logger  *slog.Logger
	service repository.ProjectRepository
}

func NewProjectService(db *gorm.DB, logger *slog.Logger, service repository.ProjectRepository) ProjectService {
	return &projectService{db: db, logger: logger, service: service}
}

func (s *projectService) Create(req *models.ProjectCreateReq) (*models.ProjectCreateResponse, error) {

	if req.Title == "" || req.Description == "" || req.Status != "" {
		s.logger.Error("передан пустой запрос",
			"op", "service.project.Create",
			"error", "empty req")
		return nil, errors.New("empty req")
	}

	statusName := map[string]bool{
		"to do":       true,
		"in progress": true,
		"done":        true,
	}

	if !statusName[req.Status] {
		s.logger.Error("invalid status",
			"op", "service.project.Create",
			"status", req.Status)
		return nil, errors.New("invalid status")
	}

	if err := s.service.CreateProject(req); err != nil {
		s.logger.Error("create failed", "err", err)
		return nil, err
	}

	inpReq := models.Project{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}

	inpRes := models.ProjectCreateResponse{
		Title:       inpReq.Title,
		Description: inpReq.Description,
		Status:      inpReq.Status,
	}

	s.logger.Info("create successful", "op", "service.project.Create", "inpRes", inpRes)
	return &inpRes, nil
}

func (s *projectService) ListProjects() ([]models.ProjectCreateResponse, error) {
	projects, err := s.service.ListProjects()

	if err != nil{
		s.logger.Error("failed get list", "err", err)
		return nil, err
	}

	return projects, nil
}

// func (s *projectService) GetProjectByID(id uint) (models.Project, error){

// }

/*CreateProject(req *models.ProjectCreateReq) error
GetProjectByID(id uint) (models.Project, error)
ListProjects() ([]models.Project, error)
UpdateProject(id uint, req models.ProjectUpdReq) error
DeleteProject(id uint) error*/
