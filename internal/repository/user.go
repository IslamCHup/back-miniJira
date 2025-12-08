package repository

import (
	"back-minijira-petproject1/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type UserRepository interface{}

type userRepository struct{
	db *gorm.DB
	logger *slog.Logger
}

func NewUserRepository(db *gorm.DB,logger *slog.Logger) UserRepository{
	return &userRepository{db:db,logger: logger}
}

func(r *userRepository) CreateUser(req models.UserCreateReq) error{
	if err := r.db.Create(&req).Error; err != nil{
		r.logger.Error("failed to create user","error",err)
	} 
	r.logger.Info("user created")
		return nil
	}


	func(r *userRepository) GetUserByID (id uint) (models.User,error){
		var user models.User

		if err := r.db.First(&user,id).Error; err != nil{
			r.logger.Error("GetUserByID failed","id",id,"error",err)
			return models.User{},err
		}
		r.logger.Info("GetUserByID","id",id)
		return user,nil
	}