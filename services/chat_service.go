package services

import (
	model "backend/models"
	"backend/utils"
)

func GetAllConversations(ID string) ([]model.ChatConversation, error) {
	var conversation []model.ChatConversation
	err := utils.GetDB().Preload("Participants"). // Preload Participants
							Preload("Messages"). // Preload Messages
							Joins("JOIN conversation_users ON chat_conversations.id = conversation_users.chat_conversation_id").
							Where("conversation_users.user_id = ?", ID).
							Find(&conversation).Error
	if err != nil {
		return nil, err
	}
	return conversation, nil
}

func GetConversationById(ID string) ([]uint64, error) {
	var userIDs []uint64
	err := utils.GetDB().Model(&model.ConversationUser{}).
		Where("chat_conversation_id = ?", ID).
		Pluck("user_id", &userIDs).Error
	if err != nil {
		return nil, err
	}
	return userIDs, nil
}
