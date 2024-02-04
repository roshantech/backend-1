package services

import (
	"encoding/json"
	"log"
	"strconv"

	model "backend/models"
	"backend/utils"

	"github.com/gofiber/websocket/v2"
)

type NotificationMessage struct {
	Type    string
	Message model.EaNotification
}

type ChatMessage struct {
	Type    string
	Message model.ChatMessage
}

func SendNotification(cNotification *model.EaNotification) error {
	err := utils.GetDB().Model(&model.EaNotification{}).Create(&cNotification).Error
	if err != nil {
		return err
	}
	payload := &NotificationMessage{
		Type:    "Notification",
		Message: *cNotification,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshaling notification:", err)
		return err
	}
	utils.Connections.Lock()
	for conn := range utils.Connections.Clients {
		val := conn.Query("ID")
		if val == strconv.Itoa(int(cNotification.UserID)) {
			err := conn.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				log.Println("Error sending notification:", err)
				continue
			}
		}
	}
	utils.Connections.Unlock()
	return nil
}

func SendMessage(message *model.ChatMessage) error {
	err := utils.GetDB().Model(&model.ChatMessage{}).Create(&message).Error
	if err != nil {
		return err
	}
	payload := &ChatMessage{
		Type:    "ChatMessage",
		Message: *message,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshaling notification:", err)
		return err
	}
	userIDs, err := GetConversationById(message.ConversationID)
	utils.Connections.Lock()
	for conn := range utils.Connections.Clients {
		val := conn.Query("ID")
		ID, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		if containsUserID(userIDs, ID) {
			err := conn.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				log.Println("Error sending notification:", err)
				continue
			}
		}
	}
	utils.Connections.Unlock()
	return nil
}

func containsUserID(userIDs []uint64, targetUserID uint64) bool {
	for _, userID := range userIDs {
		if userID == targetUserID {
			return true
		}
	}
	return false
}
