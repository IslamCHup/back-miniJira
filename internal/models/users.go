package models

type User struct {
	Base
	FullName string `json:"full_name"`
	Tasks    []Task `gorm:"many2many:user_tasks;" json:"-"`
	IsAdmin  bool   `json:"is_admin"`
}

type UserCreateReq struct {
	FullName string `json:"full_name" binding:"required"`
	TaskIDs  []uint `json:"task_ids"`
	IsAdmin  bool   `json:"is_admin" binding:"required"`
}

type UserUpdateReq struct {
	FullName *string `json:"full_name"`
	TaskIDs  []uint  `json:"task_ids"`
	IsAdmin  *bool   `json:"is_admin"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	IsAdmin  bool   `json:"is_admin"`
	TaskIDs  []uint `json:"task_ids"`
}
