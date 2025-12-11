package models

import "time"

type Task struct {
	Base
	Title        string        `json:"title" gorm:"type:varchar(255);not null"`
	Description  string        `json:"description" gorm:"type:text"`
	Status       string        `json:"status" gorm:"type:varchar(50);default:'todo';index"`
	ProjectID    uint          `json:"project_id" gorm:"index"`
	Users        []User        `json:"users" gorm:"many2many:task_users;"`
	Priority     int           `json:"priority" gorm:"default:0;index"`
	LimitUser    int           `json:"limit" gorm:"default:1"`
	StartTask    *time.Time    `json:"start_task" gorm:"index"`
	FinishTask   *time.Time    `json:"finish_task" gorm:"index"`
	ChatMessages []ChatMessage `gorm:"polymorphic:Chatable"`
}

type TaskCreateReq struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status" binding:"required,oneof=todo in_progress done"`
	ProjectID   uint       `json:"project_id"`
	// поддержка camelCase от фронта
	ProjectId   uint       `json:"projectId"`
	Users       []User     `json:"users"`
	Priority    int        `json:"priority"`
	LimitUser   int        `json:"limit"`
	StartTask   *time.Time `json:"start_task"`
	FinishTask  *time.Time `json:"finish_task"`
}

type TaskCreateRes struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ProjectID   uint   `json:"project_id"`
	Users       []User `json:"users"`
}

type TaskUpdateReq struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *string    `json:"status" binding:"omitempty,oneof=todo in_progress done"`
	Users       *[]User    `json:"users"`
	Priority    *int       `json:"priority"`
	LimitUser   *int       `json:"limit"`
	StartTask   *time.Time `json:"start_task"`
	FinishTask  *time.Time `json:"finish_task"`
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
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	ProjectID   uint       `json:"project_id"`
	Users       []User     `json:"users"`
	Priority    string     `json:"priority"`
	LimitUser   int        `json:"limit"`
	StartTask   *time.Time `json:"start_task"`
	FinishTask  *time.Time `json:"finish_task"`
}
