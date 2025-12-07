package models

import "time"

type Task struct {
	Base
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ProjectID   uint   `json:"project_id"`
	// Users       []User `json:"user" gorm:"many2many:task_users;"`
	// Comments    []Comment `json"comments"`
	StartTask  time.Time `json:"start_task"`
	FinishTask time.Time `json:"finish_task"`
}

type TaskCreateReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ProjectID   uint   `json:"project_id"`
	// Users       []User `json:"user" gorm:"many2many:task_users;"`
	// Comments    []Comment `json"comments"`
	StartTask  time.Time `json:"start_task"`
	FinishTask time.Time `json:"finish_task"`
}

type TaskCreateRes struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ProjectID   uint   `json:"project_id"`
	// Users       []User `json:"user" gorm:"many2many:task_users;"`
	// Comments    []Comment `json"comments"`
}

type TaskUpdateReq struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	ProjectID   *uint   `json:"project_id"`
	// Users       []*User `json:"user" gorm:"many2many:task_users;"`
	// Comments    []*Comment `json"comments"`
	StartTask  *time.Time `json:"start_task"`
	FinishTask *time.Time `json:"finish_task"`
}
