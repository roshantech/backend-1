package services

import (
	model "backend/models"
	"backend/utils"
)


func GetPostUsingId(Id uint) (*[]model.Post, error) {
	var post []model.Post
	err := utils.GetDB().Where("user_id = ?", Id).Find(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}


func CreatePost(post model.Post) (*model.Post, error) {
	err := utils.GetDB().Create(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}
