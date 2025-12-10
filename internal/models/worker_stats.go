package models

type WorkerStats struct {
	UserID         uint   `json:"user_id"`
	Name           string `json:"name"`
	CompletedTasks int    `json:"completed_tasks"`
}

type AvgTimeDTO struct {
	TasksCount     int    `json:"tasks_count"`
	CompletedCount int    `json:"completed_count"`
	AverageSeconds int64  `json:"average_seconds"`
	AverageHuman   string `json:"average_human"`
}

type CompletionPercentDTO struct {
	TotalTasks int     `json:"total_tasks"`
	DoneTasks  int     `json:"done_tasks"`
	Percent    float64 `json:"percent"`
}

type UserTrackerTaskDTO struct {
	TaskID    uint   `json:"task_id"`
	Title     string `json:"title"`
	StartedAt string `json:"started_at"`
}

type UserTrackerDTO struct {
	UserID             uint                 `json:"user_id"`
	InProgress         int                  `json:"in_progress"`
	Done               int                  `json:"done"`
	TotalTimeSeconds   int64                `json:"total_time_seconds"`
	TotalTimeHuman     string               `json:"total_time_human"`
	AverageTimeSeconds int64                `json:"average_time_seconds"`
	AverageTimeHuman   string               `json:"average_time_human"`
	ActiveTasks        []UserTrackerTaskDTO `json:"active_tasks"`
}
