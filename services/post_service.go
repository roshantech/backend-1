package services

import (
	model "backend/models"
	"backend/utils"

)

func CreatePost(post model.Post) (*model.Post, error) {
	err := utils.GetDB().Create(&post).Error
	if err != nil {
		
		return nil, err
	}
	return &post, nil
}
