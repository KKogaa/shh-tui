package client

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/KKogaa/shh-tui/config"
	"github.com/KKogaa/shh-tui/entity"
	"github.com/gorilla/websocket"
)

type WSClient struct {
	websocketConn *websocket.Conn
	config        *config.Config
}

func NewWsClient(config *config.Config) (WSClient, error) {

	conn, _, err := websocket.DefaultDialer.Dial(config.Server.URL, nil)
	if err != nil {
		log.Fatal("error connecting to WebSocket server:", err)
	}

	return WSClient{
		websocketConn: conn,
		config:        config,
	}, nil
}

func (w WSClient) SendMessage(message string) error {

	msg := entity.Message{
		User:     w.config.Client.Username,
		Msg:      message,
		Chatroom: "test",
	}
	jsonMsg, err := json.Marshal(msg)
	err = w.websocketConn.WriteMessage(websocket.TextMessage, []byte(jsonMsg))
	if err != nil {
		return fmt.Errorf("Error sending message to WebSocket server: %s", err)
	}

	return nil
}
func (w WSClient) ReadMessage() (entity.Message, error) {
	_, message, err := w.websocketConn.ReadMessage()
	if err != nil {
		return entity.Message{}, fmt.Errorf("error reading from websocket")
	}

	var decodedMsg entity.Message
	err = json.Unmarshal(message, &decodedMsg)
	if err != nil {
		return entity.Message{}, fmt.Errorf("error decoding json")
	}

	return decodedMsg, nil
}
