package repository

import (
	"back-minijira-petproject1/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type TeamRepository interface {
	CreateTeam(req *models.Team) error
	UpdateTeam(req *models.Team) error
	DeleteTeam(id uint) error
	GetTeamByID(teamID uint) (models.Team, []uint, error)
	AssignUsers(team *models.Team, userIDs []uint) error
	ListTeams() ([]models.Team, error)
	GetTeamsByProjectID(projectID uint) ([]models.Team, error)
}

type teamRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewTeamRepository(db *gorm.DB, logger *slog.Logger) TeamRepository {
	return &teamRepository{db: db, logger: logger}
}

func (r *teamRepository) CreateTeam(req *models.Team) error {
	res := r.db.Create(req)
	if res.Error != nil {
		r.logger.Error("Create team failed", "err", res.Error)
		return res.Error
	}
	r.logger.Info("Create team success", "team_id", req.ID)
	return nil
}

func (r *teamRepository) UpdateTeam(req *models.Team) error {
	res := r.db.Save(req)
	if res.Error != nil {
		r.logger.Error("Update team failed", "team_id", req.ID, "err", res.Error)
		return res.Error
	}
	r.logger.Info("Update team success", "team_id", req.ID, "rows", res.RowsAffected)
	return nil
}

func (r *teamRepository) DeleteTeam(id uint) error {
	res := r.db.Delete(&models.Team{}, id)
	if res.Error != nil {
		r.logger.Error("Delete team failed", "id", id, "err", res.Error)
		return res.Error
	}
	r.logger.Info("Delete team success", "id", id)
	return nil
}

func (r *teamRepository) GetTeamByID(teamID uint) (models.Team, []uint, error) {
	var team models.Team
	if err := r.db.First(&team, teamID).Error; err != nil {
		r.logger.Error("GetByIDTeam failed", "id", teamID, "err", err)
		return team, nil, err
	}

	var ids []uint
	err := r.db.Table("team_users").
		Where("team_id = ?", teamID).
		Pluck("user_id", &ids).Error

	if err != nil {
		r.logger.Error("GetUserIDs failed", "team_id", teamID, "err", err)
		return team, nil, err
	}

	r.logger.Info("GetByIDTeam success", "id", teamID)
	return team, ids, nil
}

func (r *teamRepository) AssignUsers(team *models.Team, userIDs []uint) error {
	var users []models.User

	if err := r.db.Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		r.logger.Error("AssignUsers load users failed", "err", err)
		return err
	}

	if err := r.db.Model(team).Association("Users").Replace(users); err != nil {
		r.logger.Error("AssignUsers failed", "err", err)
		return err
	}

	r.logger.Info("AssignUsers success", "team_id", team.ID)
	return nil
}

func (r *teamRepository) ListTeams() ([]models.Team, error) {
	var teams []models.Team
	if err := r.db.Find(&teams).Error; err != nil {
		r.logger.Error("ListTeams failed", "err", err)
		return nil, err
	}
	r.logger.Info("ListTeams success", "count", len(teams))
	return teams, nil
}

func (r *teamRepository) GetTeamsByProjectID(projectID uint) ([]models.Team, error) {
	var teams []models.Team
	if err := r.db.Where("project_id = ?", projectID).Find(&teams).Error; err != nil {
		r.logger.Error("GetTeamsByProjectID failed", "project_id", projectID, "err", err)
		return nil, err
	}
	r.logger.Info("GetTeamsByProjectID success", "project_id", projectID, "count", len(teams))
	return teams, nil
}
