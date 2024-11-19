package chat

import (
	"sync"

	"github.com/google/uuid"
)

type RoomType byte

const (
	Private RoomType = iota + 1
	Public
	Global
)

type Room struct {
	ID    uuid.UUID
	Type  RoomType
	Users map[*Client]bool

	Broadcast chan []byte
	Manager   *RoomManger

	mu sync.RWMutex
}

type RoomManger struct {
	Add  chan *Client
	Kick chan *Client
}

func NewRoom(roomType RoomType) *Room {
	return &Room{
		ID:        uuid.New(),
		Type:      roomType,
		Users:     map[*Client]bool{},
		Broadcast: make(chan []byte),
		Manager: &RoomManger{
			Add:  make(chan *Client),
			Kick: make(chan *Client),
		},
	}
}
