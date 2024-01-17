package utils

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

const (
	INVALID_VAL_ID  = "Invalid value for 'Id'"
	USER_NOT_FOUND  = "User not found"
	PARSE_FORM      = "Unable to parse form"
	ERR_OPEN_FILE   = "Unable to open uploaded file"
	ERR_CREATE_FILE = "Unable to create device file"
	ERR_COPY_FILE   = "Failed to copy uploaded file to device file"
)

type WebSocketConnections struct {
	Clients map[*websocket.Conn]bool
	sync.Mutex
}

var Connections = &WebSocketConnections{
	Clients: make(map[*websocket.Conn]bool),
}
