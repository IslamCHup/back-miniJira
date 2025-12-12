package models

import "time"

type Project struct {
	Base
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Tasks        *[]Task       `json:"tasks"`
	Status       string        `json:"status"`
	TimeEnd      *time.Time    `json:"time_end"`
	ChatMessages []ChatMessage `gorm:"polymorphic:Chatable"`
	Teams        []Team        `json:"teams"`
}

type ProjectCreateReq struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Tasks       *[]Task    `json:"tasks"`
	Status      string     `json:"status"`
	TimeEnd     *time.Time `json:"time_end"`
	Teams       []Team     `json:"teams"`
}

type ProjectUpdReq struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Tasks       *[]Task `json:"tasks"`
	Status      *string `json:"status"`
	Teams       *[]Team `json:"teams"`
}

type ProjectCreateResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	TimeEnd     time.Time `json:"time_end"`
	Teams       []Team    `json:"teams"`
}

type ProjectFilter struct {
	Title       *string
	Description *string
	Status      *string
	Limit       int
	Offset      int
}
