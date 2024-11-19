package chat

import "github.com/google/uuid"

type Client struct {
	ID       uuid.UUID
	Username string

	ChatRoom *Room
	recieve  chan []byte
}

func NewClient(username string, chatRoom *Room) *Client {
	return &Client{
		ID:       uuid.New(),
		Username: username,
		ChatRoom: chatRoom,
		recieve:  make(chan []byte),
	}
}
