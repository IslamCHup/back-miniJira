package service

import (
	"back-minijira-petproject1/internal/models"
	"back-minijira-petproject1/internal/repository"
	"fmt"
	"log/slog"
	"time"
)

type ReportService interface {
	TopWorkers(projectID uint) ([]models.WorkerStats, error)
	AverageTime(projectID uint) (models.AvgTimeDTO, error)
	CompletionPercent(projectID uint) (models.CompletionPercentDTO, error)
	UserTracker(projectID uint, userID uint) (models.UserTrackerDTO, error)
}

type reportService struct {
	repo repository.ReportRepository
	logger slog.Logger
}

func NewReportService(report repository.ReportRepository, logger *slog.Logger) ReportService {
	return &reportService{repo: report}
}

func (s *reportService) TopWorkers(projectID uint) ([]models.WorkerStats, error) {
	return s.repo.GetTopWorkers(projectID)
}

func (s *reportService) AverageTime(projectID uint) (models.AvgTimeDTO, error) {
	tasks, err := s.repo.GetCompletedTasksTimes(projectID)
	if err != nil {
		return models.AvgTimeDTO{}, err
	}

	if len(tasks) == 0 {
		return models.AvgTimeDTO{
			TasksCount: 0, CompletedCount: 0, AverageSeconds: 0, AverageHuman: "0s",
		}, nil
	}

	var total time.Duration
	for _, t := range tasks {
		total += t.FinishTask.Sub(*t.StartTask)
	}
	avg := total / time.Duration(len(tasks))

	return models.AvgTimeDTO{
		TasksCount:     len(tasks),
		CompletedCount: len(tasks),
		AverageSeconds: int64(avg.Seconds()),
		AverageHuman:   avg.Round(time.Second).String(),
	}, nil
}

// --- Completion Percent ---
func (s *reportService) CompletionPercent(projectID uint) (models.CompletionPercentDTO, error) {
	total, err := s.repo.CountTasks(projectID)
	if err != nil {
		return models.CompletionPercentDTO{}, err
	}
	done, err := s.repo.CountDoneTasks(projectID)
	if err != nil {
		return models.CompletionPercentDTO{}, err
	}

	percent := 0.0
	if total > 0 {
		percent = float64(done) / float64(total) * 100
	}

	return models.CompletionPercentDTO{TotalTasks: total, DoneTasks: done, Percent: percent}, nil
}

// --- User Tracker ---
func (s *reportService) UserTracker(projectID uint, userID uint) (models.UserTrackerDTO, error) {
	tasks, err := s.repo.GetUserTasks(projectID, userID)
	if err != nil {
		return models.UserTrackerDTO{}, err
	}

	tracker := models.UserTrackerDTO{UserID: userID}
	var sum int64
	var doneCount int64

	for _, t := range tasks {
		if t.Status == "in_progress" {
			tracker.InProgress++
			tracker.ActiveTasks = append(tracker.ActiveTasks, models.UserTrackerTaskDTO{
				TaskID:    t.ID,
				Title:     t.Title,
				StartedAt: t.StartTask.Format(time.RFC3339),
			})
		}

		if t.Status == "done" {
			tracker.Done++
			doneCount++
			sum += int64(t.FinishTask.Sub(*t.StartTask).Seconds())
		}
	}

	tracker.TotalTimeSeconds = sum
	tracker.TotalTimeHuman = fmt.Sprintf("%ds", sum)

	if doneCount > 0 {
		avg := sum / doneCount
		tracker.AverageTimeSeconds = avg
		tracker.AverageTimeHuman = fmt.Sprintf("%ds", avg)
	}

	return tracker, nil
}
