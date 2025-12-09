package models

import "time"

type Task struct {
	Base
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	Status       string         `json:"status" gorm:"default:to do"`
	ProjectID    uint           `json:"project_id"`
	Users        []UserResponse `json:"user" gorm:"many2many:task_users;"`
	Priority     int            `json:"priority" gorm:"default:0;"`
	LimitUser    int            `json:"limit"`
	StartTask    time.Time      `json:"start_task"`
	FinishTask   time.Time      `json:"finish_task"`
	ChatMessages []ChatMessage  `gorm:"polymorphic:Chatable"`
}

type TaskCreateReq struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Status      string         `json:"status" gorm:"default:to do"`
	ProjectID   uint           `json:"project_id"`
	Users       []UserResponse `json:"user" gorm:"many2many:task_users;"`
	Priority    int            `json:"priority" gorm:"default:0;"`
	LimitUser   int            `json:"limit"`
	StartTask   time.Time      `json:"start_task"`
	FinishTask  time.Time      `json:"finish_task"`
}

type TaskCreateRes struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	ProjectID   uint           `json:"project_id"`
	Users       []UserResponse `json:"user"`
}

type TaskUpdateReq struct {
	Title       *string         `json:"title"`
	Description *string         `json:"description"`
	Status      *string         `json:"status"`
	ProjectID   *uint           `json:"project_id"`
	Users       []*UserResponse `json:"user"`
	Priority    *int            `json:"priority"`
	LimitUser   *int            `json:"limit"`
	StartTask   *time.Time      `json:"start_task"`
	FinishTask  *time.Time      `json:"finish_task"`
}

type TaskFilter struct {
	Status    *string
	UserID    *uint
	ProjectID *uint
	Search    *string
	Priority  *int
	SortBy    *string
	SortOrder *string
	Limit     int
	Offset    int
}

type TaskResponse struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	ProjectID   uint           `json:"project_id"`
	Users       []UserResponse `json:"user"`
	Priority    string         `json:"priority"`
	LimitUser   int            `json:"limit"`
	StartTask   time.Time      `json:"start_task"`
	FinishTask  time.Time      `json:"finish_task"`
}
