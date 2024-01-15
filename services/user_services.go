package services

import (
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
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return &user, nil
}

func CreateUser(user model.User) (*model.User, error) {
	hashedPassword, err := model.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	err = utils.GetDB().Table("users").Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func ValidateEmailId(email string) (bool, error) {

	_, err := GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := utils.GetDB().Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(ID string) (*model.User, error) {
	var user model.User
	err := utils.GetDB().Where("ID = ?", ID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUserById(user model.User) (*model.User, error) {
	err := utils.GetDB().Where("ID = ?", user.ID).Updates(&user).Error
	if err != nil {
		return nil, err
	}
	address := user.Address
	err = utils.GetDB().Where("user_id = ?", user.ID).Updates(&address).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
