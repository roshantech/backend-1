package services

import (
	model "backend/models"
	"backend/utils"
	"errors"

	"gorm.io/gorm"
)

func GetPostUsingUserId(Id uint) (*[]model.Post, error) {
	var post []model.Post
	err := utils.GetDB().Preload("Likes").Preload("Comments").Find(&post).Error
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

func GetPostUsingId(Id uint) (*model.Post, error) {
	var post model.Post
	err := utils.GetDB().Where("ID = ?", Id).Find(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func UpdaatePostLikes(postID uint, userID uint) error {
	var existingLike model.Like
	var err error
	if err = utils.GetDB().Where("post_id = ? AND user_id = ?", postID, userID).First(&existingLike).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newLike := model.Like{PostID: postID, UserID: userID}
		if err = utils.GetDB().Create(&newLike).Error; err != nil {
			return err
		}
	} else if err == nil {
		if err := utils.GetDB().Delete(&existingLike).Error; err != nil {
			return err
		}

	}

	return nil
}


func CreatePostComments(comment model.Comment) error {

	if err := utils.GetDB().Create(&comment).Error; err != nil {
			return err
	}

	return nil
}