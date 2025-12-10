package models

type Team struct {
	Base
	Name      string `json:"name"`
	Users     []User `gorm:"many2many:team_users;" json:"-"`
	ProjectID uint   `json:"project_id"`
	UserID    uint   `json:"user_id"`
}

type TeamCreateReq struct {
	Name      string `json:"name" binding:"required"`
	Users     []uint `json:"user_ids" binding:"required"`
	ProjectID uint   `json:"project_id"`
	UserID    uint   `json:"user_id"`
}

type TeamUpdateReq struct {
	Name      *string `json:"name"`
	Users     *[]uint `json:"user_ids"`
	ProjectID *uint   `json:"project_id"`
	UserID    *uint   `json:"user_id"`
}
