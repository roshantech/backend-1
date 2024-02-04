package services

import (
	model "backend/models"
	"backend/utils"
	"strconv"
	"time"

	"gorm.io/gorm"
	//"backend/utils"
)

type Follower struct {
	gorm.Model
	UserID     uint       `json:"user_id"`
	FollowerID uint       `json:"follower_id"`
	user       model.User `json:"user"`
}

func GetUserFollowers(userID uint) (*[]Follower, error) {
	var follower []Follower
	if err := utils.GetDB().Model(model.Follower{}).Where("user_id = ?", userID).Find(&follower).Error; err != nil {
		return nil, err
	}

	for key, _ := range follower {
		if err := utils.GetDB().Model(model.User{}).Where("id = ?", follower[key].FollowerID).Find(&follower[key].user).Error; err != nil {
			return nil, err
		}
	}
	return &follower, nil
}

func AddFollowerToUser(userID uint, followerID uint) error {
	following := &model.Following{
		UserID:      userID,
		FollowingID: followerID,
	}
	follower := &model.Following{
		UserID:      followerID,
		FollowingID: userID,
	}
	if err := utils.GetDB().Create(following).Error; err != nil {
		return err
	}
	if err := utils.GetDB().Create(follower).Error; err != nil {
		return err
	}

	user, err := GetUserByID(strconv.Itoa(int(userID)))
	if err != nil {
		return err
	}
	notifyUser := &model.EaNotification{
		UserID:          followerID,
		Title:           user.Username,
		NotifyMsg:       " Is Following You",
		Status:          "Unread",
		RecordTimestamp: time.Now().Format("2006-01-02 15:04:05"),
	}
	err = SendNotification(notifyUser)
	if err != nil {
		return err
	}

	return nil
}
