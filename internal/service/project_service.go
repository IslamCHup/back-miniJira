package service

import (
	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/repository"
	"errors"
	"log/slog"
	"slices"
	"strings"

	"gorm.io/gorm"
)

type ProjectService interface {
	Create(req *models.ProjectCreateReq) (*models.ProjectCreateResponse, error)
	ListProjects(filter *models.ProjectFilter) ([]models.ProjectCreateResponse, error)
	GetByID(id uint) (models.ProjectCreateResponse, error)
	Delete(id uint) error
	UpdateProject(id uint, req models.ProjectUpdReq) error
}

type projectService struct {
	db     *gorm.DB
	logger *slog.Logger
	repo   repository.ProjectRepository
}

func NewProjectService(db *gorm.DB, logger *slog.Logger, repo repository.ProjectRepository) ProjectService {
	return &projectService{db: db, logger: logger, repo: repo}
}

func (s *projectService) Create(req *models.ProjectCreateReq) (*models.ProjectCreateResponse, error) {

	if req.Title == "" || req.Description == "" || req.Status == "" {
		s.logger.Error("передан пустой запрос",
			"op", "service.project.Create",
			"error", "empty req")
		return nil, errors.New("empty req")
	}

	statusName := map[string]bool{
		"todo":        true,
		"in_progress": true,
		"done":        true,
	}

	if !statusName[req.Status] {
		s.logger.Error("invalid status",
			"op", "service.project.Create",
			"status", req.Status)
		return nil, errors.New("invalid status")
	}

	if err := s.repo.CreateProject(req); err != nil {
		s.logger.Error("create failed", "err", err)
		return nil, err
	}

	resp := models.ProjectCreateResponse{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}
	if req.TimeEnd != nil {
		resp.TimeEnd = *req.TimeEnd
	}

	s.logger.Info("create successful", "op", "service.project.Create", "resp", resp)
	return &resp, nil
}

func (s *projectService) ListProjects(filter *models.ProjectFilter) ([]models.ProjectCreateResponse, error) {
	projects, err := s.repo.ListProjects(filter)

	if err != nil {
		s.logger.Error("failed get list", "err", err)
		return nil, err
	}

	s.logger.Info("get list successful", "op", "service.project.ListProjects", "count", len(projects))
	return projects, nil
}

func (s *projectService) GetByID(id uint) (models.ProjectCreateResponse, error) {
	project, err := s.repo.GetProjectByID(id)

	if err != nil {
		s.logger.Error("failed get project by id", "id", id, "err", err)
		return models.ProjectCreateResponse{}, err
	}

	s.logger.Info("get project by id successful", "op", "service.project.GetByID", "id", id, "project", project)
	return project, nil
}

func (s *projectService) Delete(id uint) error {
	if err := s.repo.DeleteProject(id); err != nil {
		s.logger.Error("failed delete by id", "id", id, "err", err)
		return err
	}
	s.logger.Info("delete project by id successful", "op", "service.project.deleteProject")
	return nil
}

func (s *projectService) UpdateProject(id uint, req models.ProjectUpdReq) error {
	project, err := s.repo.GetProjectByID(id)
	if err != nil {
		s.logger.Error("failed get project by id", "id", id, "err", err)
		return err
	}

	allowedStatus := map[string][]string{
		"todo":        {"in_progress"},
		"in_progress": {"todo", "done"},
		"done":        {"in_progress"},
	}
	if req.Status != nil {
		oldStatusTask := strings.ToLower(project.Status)
		var newStatusProject string
		if req.Status != nil {
			newStatusProject = strings.ToLower(*req.Status)
			statusSlice, ok := allowedStatus[oldStatusTask]
			if !ok {
				s.logger.Error("such status does not exist", "task_status_current", project.Status)
				return errors.New("the task status changes only in a certain order")
			}

			if !slices.Contains(statusSlice, newStatusProject) {
				s.logger.Error("can't skip status", "task_status_current", project.Status, "req_tas_status", req.Status)
				return errors.New("the task status changes only in a certain order")
			}
		}
	}
	if err := s.repo.UpdateProject(id, req); err != nil {
		s.logger.Error("failed update project", "err", err)
		return err
	}

	s.logger.Info("update project by id successful", "op", "service.project.updateProject", "id", id)
	return nil
}
