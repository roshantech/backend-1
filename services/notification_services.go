package services

import (
	"encoding/json"
	"log"

	model "backend/models"
	"backend/utils"

	"github.com/gofiber/websocket/v2"
)

func SendNotification(cNotification *model.EaNotification) error {
	err := utils.GetDB().Model(&model.EaNotification{}).Create(&cNotification).Error
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(cNotification)
	if err != nil {
		log.Println("Error marshaling notification:", err)
		return err
	}
	utils.Connections.Lock()
	for conn := range utils.Connections.Clients {
		err := conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Println("Error sending notification:", err)
			continue
		}
	}
	utils.Connections.Unlock()
	return nil
}
