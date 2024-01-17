package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"backend/dto"
	model "backend/models"
	"backend/utils"
	"backend/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func NotificationWs(c *fiber.Ctx) error {
	return websocket.New(func(conn *websocket.Conn) {
		// Add the WebSocket connection to the connections map
		var contains bool
		for con := range utils.Connections.Clients {
			if conn.LocalAddr() == con.LocalAddr() {
				contains = true
			}
		}
		if !contains {
			utils.Connections.Lock()
			utils.Connections.Clients[conn] = true
			utils.Connections.Unlock()
		}

		defer func() {
			utils.Connections.Lock()
			delete(utils.Connections.Clients, conn)
			utils.Connections.Unlock()
		}()

		for {
			_, P, err := conn.ReadMessage()
			if err != nil {
				break
			}
			var readAll []int
			err = json.Unmarshal(P, &readAll)
			if err != nil {
				break
			}
			err = markNotificcationAsRead(readAll)
			if err != nil {
				break
			}
		}
	})(c)
}

func saveMessage(notification *dto.Notificationdto) error {

	return nil
}

func markNotificcationAsRead(ID []int) error {

	for _, id := range ID {
		err := utils.GetDB().Model(model.EaNotification{}).Where("id = ?", id).Update("status", "Read").Error
		if err != nil {
			return err
		}
	}
	return nil
}

func GetAllNotifications(c *fiber.Ctx) error {
	idStr := c.Query("Id")

	id, err := strconv.ParseUint(idStr, 10, 16)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid value for 'Id' parameter")
	}
	var notifications []model.EaNotification
	err = utils.GetDB().Model(model.EaNotification{}).Where("user_id = ?", id).Scan(&notifications).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error Fetching Data"})
	}
	return c.JSON(notifications)
}

func SendBrodcastNotification(c *fiber.Ctx) error {
	var brodcast struct {
		Message      string
		BrodcastType string
	}
	err := c.BodyParser(&brodcast)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Unable to parse Data")
	}
	var user []model.User
	err = utils.GetDB().Model([]model.User{}).Scan(&user).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Unable To get UserInfo")
	}

	
	return c.SendString("Broadcast Sent Successfully")
}

func SendNotification(c *fiber.Ctx) error {
	var notification dto.Notificationdto
	err := c.BodyParser(&notification)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Unable to parse Data")
	}
	cNotification := &model.EaNotification{
		UserID:          notification.UserID,
		NotifyMsg:       notification.Message,
		Status:          "Unread",
		RecordTimestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	err = services.SendNotification(cNotification)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Unable to send Notification : " + err.Error())
	}

	return nil
}
