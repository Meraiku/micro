package chat

import (
	"context"
	"embed"
	"sync"

	"github.com/google/uuid"
	"github.com/meraiku/micro/pkg/logging"
)

//go:embed templates/*.html
var templates embed.FS

type MessageRepository interface {
	Save(roomID uuid.UUID, message *Message) error
	GetAll(roomID uuid.UUID) ([]Message, error)
}

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

	Broadcast chan *Message
	Manager   *RoomManger

	msgRepo MessageRepository
	mu      sync.RWMutex
}

type RoomManger struct {
	Add    chan *Client
	Logout chan *Client
	Kick   chan *Client
}

func NewRoom(roomType RoomType, msgRepo MessageRepository) *Room {
	return &Room{
		ID:        uuid.New(),
		Type:      roomType,
		Users:     map[*Client]bool{},
		Broadcast: make(chan *Message, 100),
		Manager: &RoomManger{
			Add:    make(chan *Client),
			Logout: make(chan *Client),
			Kick:   make(chan *Client),
		},
	}
}

func (r *Room) Run(ctx context.Context) {
	logging.WithAttrs(
		ctx,
		logging.String("room_id", r.ID.String()),
	)

	log := logging.L(ctx)

	for {
		select {
		case client := <-r.Manager.Add:
			r.Add(client)
		case client := <-r.Manager.Logout:
			r.Logout(client)
		case client := <-r.Manager.Kick:
			r.Kick(client)
		case msg := <-r.Broadcast:

			for c, ok := range r.Users {
				if ok {

					log.Debug(
						"sending message",
						logging.String("client_id", c.ID.String()),
						logging.Any("message", msg),
					)
					c.recieve <- msg.Render()
					continue
				}

				log.Info(
					"client left room",
					logging.String("client_id", c.ID.String()),
				)

				delete(r.Users, c)
				close(c.recieve)
			}
		}

	}
}

func (r *Room) Add(client *Client) {
	r.mu.Lock()
	r.Users[client] = true
	r.mu.Unlock()
}

func (r *Room) Logout(client *Client) {
	r.mu.Lock()
	r.Users[client] = false
	r.mu.Unlock()
}

func (r *Room) Kick(client *Client) {
	r.mu.Lock()
	delete(r.Users, client)
	r.mu.Unlock()
}
