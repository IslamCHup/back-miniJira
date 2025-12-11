package service

import (
	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/repository"
	"errors"
	"fmt"
	"log/slog"
)

type TeamService interface {
	Create(req models.TeamCreateReq, currentUser models.User) (*models.TeamResponse, error)
	Update(id uint, req models.TeamUpdateReq, currentUser models.User) error
	Delete(id uint, currentUser models.User) error
	GetByID(id uint) (*models.TeamResponse, error)
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

func (s *teamService) checkPermission(currentUser models.User, team models.Team) error {
	if currentUser.IsAdmin || currentUser.ID == team.UserID {
		return nil
	}

	s.logger.Warn("team permission denied",
		"current_user_id", currentUser.ID,
		"team_lead_id", team.UserID,
		"team_id", team.ID,
	)

	return errors.New("недостаточно прав")
}

func (s *teamService) Create(req models.TeamCreateReq, currentUser models.User) (*models.TeamResponse, error) {
	if err := validateTeamName(req.Name); err != nil {
		s.logger.Error("CreateTeam validate name failed", "name", req.Name, "err", err)
		return nil, err
	}

	if !currentUser.IsAdmin && currentUser.ID != req.UserID {
		err := errors.New("недостаточно прав")
		s.logger.Warn("CreateTeam permission denied",
			"current_user_id", currentUser.ID,
			"lead_id", req.UserID,
		)
		return nil, err
	}

	team := &models.Team{
		Name:      req.Name,
		ProjectID: req.ProjectID,
		UserID:    req.UserID,
	}

	if err := s.repo.CreateTeam(team); err != nil {
		s.logger.Error("CreateTeam failed", "err", err)
		return nil, err
	}

	if err := s.repo.AssignUsers(team, req.Users); err != nil {
		s.logger.Error("CreateTeam assign users failed", "team_id", team.ID, "err", err)
		return nil, err
	}

	savedTeam, userIDs, err := s.repo.GetTeamByID(team.ID)
	if err != nil {
		s.logger.Error("CreateTeam reload team failed", "team_id", team.ID, "err", err)
		return nil, err
	}

	resp := models.TeamResponse{
		ID:        savedTeam.ID,
		Name:      savedTeam.Name,
		ProjectID: savedTeam.ProjectID,
		UserID:    savedTeam.UserID,
		UserIDs:   userIDs,
	}

	s.logger.Info("CreateTeam success", "team_id", resp.ID)
	return &resp, nil
}

func (s *teamService) Update(id uint, req models.TeamUpdateReq, currentUser models.User) error {
	team, _, err := s.repo.GetTeamByID(id)
	if err != nil {
		s.logger.Error("UpdateTeam GetTeamByID failed", "team_id", id, "err", err)
		return err
	}

	if err := s.checkPermission(currentUser, team); err != nil {

		return err
	}

	if req.Name != nil {
		if err := validateTeamName(*req.Name); err != nil {
			s.logger.Error("UpdateTeam validate name failed", "team_id", id, "name", *req.Name, "err", err)
			return err
		}
		team.Name = *req.Name
	}

	if req.ProjectID != nil {
		team.ProjectID = *req.ProjectID
	}

	if req.Users != nil {
		if err := s.repo.AssignUsers(&team, *req.Users); err != nil {
			s.logger.Error("UpdateTeam assign users failed", "team_id", id, "err", err)
			return err
		}
	}

	if req.UserID != nil {
		if !currentUser.IsAdmin {
			err := errors.New("недостаточно прав")
			s.logger.Warn("UpdateTeam change lead denied",
				"team_id", id,
				"current_user_id", currentUser.ID,
				"new_lead_id", *req.UserID,
			)
			return err
		}
		team.UserID = *req.UserID
	}

	if err := s.repo.UpdateTeam(&team); err != nil {
		s.logger.Error("UpdateTeam failed", "team_id", id, "err", err)
		return err
	}

	s.logger.Info("UpdateTeam success", "team_id", id)
	return nil
}

func (s *teamService) Delete(id uint, currentUser models.User) error {
	team, _, err := s.repo.GetTeamByID(id)
	if err != nil {
		s.logger.Error("DeleteTeam GetTeamByID failed", "team_id", id, "err", err)
		return err
	}

	if err := s.checkPermission(currentUser, team); err != nil {
		return err
	}

	if err := s.repo.DeleteTeam(id); err != nil {
		s.logger.Error("DeleteTeam failed", "team_id", id, "err", err)
		return err
	}

	s.logger.Info("DeleteTeam success", "team_id", id)
	return nil
}

func (s *teamService) GetByID(id uint) (*models.TeamResponse, error) {
	team, userIDs, err := s.repo.GetTeamByID(id)
	if err != nil {
		s.logger.Error("GetTeamByID failed", "team_id", id, "err", err)
		return nil, err
	}

	resp := &models.TeamResponse{
		ID:        team.ID,
		Name:      team.Name,
		ProjectID: team.ProjectID,
		UserID:    team.UserID,
		UserIDs:   userIDs,
	}

	s.logger.Info("GetTeamByID success", "team_id", id)
	return resp, nil
}
