package repository

import (
	"back-minijira-petproject1/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type ReportRepository interface {
	GetTopWorkers(projectID uint) ([]models.WorkerStats, error)
	GetCompletedTasksTimes(projectID uint) ([]models.Task, error)
	CountTasks(projectID uint) (int, error)
	CountDoneTasks(projectID uint) (int, error)
	GetUserTasks(projectID uint, userID uint) ([]models.Task, error)
}

type reportRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewReportRepository(db *gorm.DB, logger *slog.Logger) ReportRepository {
	return &reportRepo{db: db, logger: logger}
}

func (r *reportRepo) GetTopWorkers(projectID uint) ([]models.WorkerStats, error) {
	var result []models.WorkerStats

	err := r.db.Model(&models.Task{}).
		Select("users.id as user_id, users.name as name, COUNT(tasks.id) as completed_tasks").
		Joins("JOIN task_users tu ON tu.task_id = tasks.id").
		Joins("JOIN users ON users.id = tu.user_id").
		Where("tasks.project_id = ? AND tasks.status = 'done'", projectID).
		Group("users.id, users.name").
		Order("completed_tasks DESC").
		Scan(&result).Error

	return result, err
}

func (r *reportRepo) GetCompletedTasksTimes(projectID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.
		Where(
			"project_id = ? AND status = ? AND start_task IS NOT NULL AND finish_task IS NOT NULL",
			projectID, "done").
		Find(&tasks).Error
	return tasks, err
}

func (r *reportRepo) CountTasks(projectID uint) (int, error) {
	var count int64
	err := r.db.Model(&models.Task{}).Where("project_id = ?", projectID).Count(&count).Error
	return int(count), err
}

func (r *reportRepo) CountDoneTasks(projectID uint) (int, error) {
	var count int64
	err := r.db.Model(&models.Task{}).Where("project_id = ? AND status = ?", projectID, "done").Count(&count).Error
	return int(count), err
}

func (r *reportRepo) GetUserTasks(projectID uint, userID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Joins("JOIN task_users  ON task_users.task_id = tasks.id").
		Where("tasks.project_id = ? AND task_users.user_id = ?", projectID, userID).
		Preload("Users").
		Find(&tasks).Error
	return tasks, err
}
