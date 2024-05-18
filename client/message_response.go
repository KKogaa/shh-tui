package client

import (
	"time"

	"github.com/KKogaa/shh-tui/entity"
)

type MessageResponse struct {
	ChatroomId int64     `json:"chatroom_id"`
	CreatedAt  time.Time `json:"created_at"`
	Id         int64     `json:"id"`
	Payload    string    `json:"payload"`
	Username   string    `json:"username"`
}

func ToMessage(message MessageResponse) entity.Message {
	return entity.Message{
		Msg:      message.Payload,
		User:     message.Username,
		Chatroom: string(message.ChatroomId),
	}
}
