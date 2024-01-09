package services

import (
	"backend/dto"
	model "backend/models"
	"backend/utils"
	"errors"

	"gorm.io/gorm"
)



func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := utils.GetDB().Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return  nil ,errors.New("User not found")
		}
		return nil, err
	}
	return &user ,nil
}


func CreateUser(user dto.UserCreateUpdate) (*model.User, error) {
	hashedPassword,err  := model.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword;
	var User model.User 
	err = utils.GetDB().Table("users").Create(&User).Error
	if err != nil {
		return nil,err
	}

	return &User, nil
}