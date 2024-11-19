package chat

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/meraiku/micro/pkg/logging"
)

var (
	ErrRoomNotFound       = errors.New("room not found")
	ErrRoomClosed         = errors.New("room closed")
	ErrClientNotAvailable = errors.New("client not available")
)

var upgrader = websocket.Upgrader{

	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Service struct {
	Rooms   map[uuid.UUID]*Room
	Global  *Room
	msgRepo MessageRepository
	mu      sync.RWMutex
}

func NewService(ctx context.Context, msgRepo MessageRepository) *Service {
	global := NewRoom(Global, msgRepo)

	rooms := make(map[uuid.UUID]*Room)
	rooms[global.ID] = global

	go global.Run(ctx)

	return &Service{
		Rooms:   rooms,
		Global:  global,
		msgRepo: msgRepo,
		mu:      sync.RWMutex{},
	}
}

func (s *Service) CreateRoom(roomType RoomType) uuid.UUID {
	s.mu.Lock()
	defer s.mu.Unlock()

	room := NewRoom(roomType, s.msgRepo)
	s.Rooms[room.ID] = room

	return room.ID
}

func (s *Service) ConnectGlobal(
	client *Client,
	w http.ResponseWriter,
	r *http.Request,
) error {
	log := logging.L(r.Context())

	log.Debug(
		"upgrading connection",
		logging.String("username", client.Username),
	)
	if err := s.connect(client, w, r); err != nil {
		return err
	}

	log.Debug(
		"adding client to global room",
		logging.String("username", client.Username),
	)

	room, ok := s.Rooms[s.Global.ID]
	if !ok {
		log.Warn(
			"global room not found",
			logging.String("username", client.Username),
			logging.String("room_id", s.Global.ID.String()),
		)
		return ErrRoomNotFound
	}

	room.Add(client)
	client.ChatRoom = room

	log.Debug(
		"starting session",
		logging.String("username", client.Username),
	)

	return client.StartSession(r.Context(), client.conn)
}

func (s *Service) connect(
	client *Client,
	w http.ResponseWriter,
	r *http.Request,
) error {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("upgrading connection: %v", err)
	}

	client.addConnection(conn)

	return nil
}
