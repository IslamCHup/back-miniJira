package models

type Task struct {
	Base
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ProjectID   uint   `json:"project_id"`
	Users       []User `json:"user" gorm:"many2many:task_users;"`
}
