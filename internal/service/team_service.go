package service

import (
	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/repository"
	"errors"
	"fmt"
	"log/slog"
)

type TeamService interface {
	CreateTeam(req models.TeamCreateReq, currentUser models.User) (*models.TeamResponse, error)
	UpdateTeam(id uint, req models.TeamUpdateReq, currentUser models.User) error
	DeleteTeam(id uint, currentUser models.User) error
	GetTeamByID(id uint) (*models.TeamResponse, error)
}

type teamService struct {
	repo   repository.TeamRepository
	logger *slog.Logger
}

func NewTeamService(repo repository.TeamRepository, logger *slog.Logger) TeamService {
	return &teamService{repo: repo, logger: logger}
}

func validateTeamName(name string) error {
	switch name {
	case "testers", "backend", "frontend":
		return nil
	default:
		return fmt.Errorf("invalid team name '%s': allowed testers, backend, frontend", name)
	}
}

func (s *teamService) checkPermission(user models.User, team models.Team) error {
	if user.IsAdmin {
		return nil
	}
	if user.ID == team.UserID {
		return nil
	}
	return errors.New("недостаточно прав")
}

func (s *teamService) CreateTeam(req models.TeamCreateReq, currentUser models.User) (*models.TeamResponse, error) {
	if err := validateTeamName(req.Name); err != nil {
		return nil, err
	}

	if !currentUser.IsAdmin && currentUser.ID != req.UserID {
		return nil, errors.New("недостаточно прав")
	}

	team := &models.Team{
		Name:      req.Name,
		ProjectID: req.ProjectID,
		UserID:    req.UserID,
	}

	if err := s.repo.CreateTeam(team); err != nil {
		return nil, err
	}

	if err := s.repo.AssignUsers(team, req.Users); err != nil {
		return nil, err
	}

	savedTeam, userIDs, err := s.repo.GetTeamByID(team.ID)
	if err != nil {
		return nil, err
	}

	resp := models.TeamResponse{
		ID:        savedTeam.ID,
		Name:      savedTeam.Name,
		ProjectID: savedTeam.ProjectID,
		UserID:    savedTeam.UserID,
		UserIDs:   userIDs,
	}

	return &resp, nil
}

func (s *teamService) UpdateTeam(id uint, req models.TeamUpdateReq, currentUser models.User) error {
	team, _, err := s.repo.GetTeamByID(id)
	if err != nil {
		return err
	}

	if err := s.checkPermission(currentUser, team); err != nil {
		return err
	}

	if req.Name != nil {
		if err := validateTeamName(*req.Name); err != nil {
			return err
		}
		team.Name = *req.Name
	}

	if req.ProjectID != nil {
		team.ProjectID = *req.ProjectID
	}

	if req.Users != nil {
		if err := s.repo.AssignUsers(&team, *req.Users); err != nil {
			return err
		}
	}

	if req.UserID != nil {
		if !currentUser.IsAdmin {
			return errors.New("недостаточно прав")
		}
		team.UserID = *req.UserID
	}

	return s.repo.UpdateTeam(&team)
}

func (s *teamService) DeleteTeam(id uint, currentUser models.User) error {
	team, _, err := s.repo.GetTeamByID(id)
	if err != nil {
		return err
	}

	if err := s.checkPermission(currentUser, team); err != nil {
		return err
	}

	return s.repo.DeleteTeam(id)
}

func (s *teamService) GetTeamByID(id uint) (*models.TeamResponse, error) {
	team, userIDs, err := s.repo.GetTeamByID(id)
	if err != nil {
		return nil, err
	}

	return &models.TeamResponse{
		ID:        team.ID,
		Name:      team.Name,
		ProjectID: team.ProjectID,
		UserID:    team.UserID,
		UserIDs:   userIDs,
	}, nil
}
