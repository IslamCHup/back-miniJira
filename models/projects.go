package models

import "time"

type Project struct {
	Base
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tasks       []Task    `json:"tasks"`
	Status      string    `json:"status"`
	TimeEnd     time.Time `json:"time_end"`
	ChatMessages []ChatMessage `gorm:"polymorphic:Chatable"`
	// Comments    []Comment `json"comments"`
	// Commands []Command `json:"command"`
}

type ProjectCreateReq struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tasks       []*Task    `json:"tasks"`
	Status      string    `json:"status"`
	TimeEnd     time.Time `json:"time_end"`
	// Commands []Command `json:"command"`
}

type ProjectUpdReq struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Tasks       []*Task    `json:"tasks"`
	Status      *string    `json:"status"`
	TimeEnd     *time.Time `json:"time_end"`
	// Comments    []Comment `json"comments"`
	// Commands []Command `json:"command"`
}

type ProjectCreateResponse struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	TimeEnd     *time.Time `json:"time_end"`
	// Comments    []Comment `json"comments"`
	// Commands []Command `json:"command"`
}
