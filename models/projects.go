package models

import "time"

type Project struct {
	Base
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tasks       []Task    `json:"tasks"`
	Status      string    `json:"status"`
	TimeEnd     time.Time `json:"time_end"`
}

type ProjectCreateReq struct{
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tasks       []Task    `json:"tasks"`
	Status      string    `json:"status"`
	TimeEnd     time.Time `json:"time_end"`
}

type ProjectUpdReq struct{
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Tasks       []*Task    `json:"tasks"`
	Status      *string    `json:"status"`
	TimeEnd     *time.Time `json:"time_end"`
}

type ProjectCreateResponse struct{
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	TimeEnd     *time.Time `json:"time_end"`
}