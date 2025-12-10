package models

type User struct {
	Base
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Tasks        []Task `gorm:"many2many:user_tasks;" json:"-"`
	IsAdmin      bool   `json:"is_admin"`
	IsVerified   bool   `json:"is_verified"`
	VerifyToken  string `json:"-"`
}

type UserCreateReq struct {
	FullName string `json:"full_name" binding:"required"`
	TaskIDs  []uint `json:"task_ids"`
}

type UserUpdateReq struct {
	FullName *string `json:"full_name"`
	TaskIDs  []uint  `json:"task_ids"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
	TaskIDs  []uint `json:"task_ids"`
}
