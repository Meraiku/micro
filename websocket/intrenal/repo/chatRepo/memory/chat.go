package memory

import (
	"sync"

	"github.com/google/uuid"
	chatrepo "github.com/meraiku/micro/websocket/intrenal/repo/chatRepo"
	"github.com/meraiku/micro/websocket/intrenal/services/chat"
)

type Repository struct {
	mu       sync.RWMutex
	messages map[uuid.UUID][]chat.Message
}

func NewRepository() *Repository {
	return &Repository{
		messages: map[uuid.UUID][]chat.Message{},
	}
}

func (r *Repository) Save(roomID uuid.UUID, message *chat.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	msgs, ok := r.messages[roomID]
	if !ok {
		msgs = []chat.Message{}
	}
	msgs = append(msgs, *message)
	r.messages[roomID] = msgs

	return nil
}

func (r *Repository) GetAll(roomID uuid.UUID) ([]chat.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	msgs, ok := r.messages[roomID]
	if !ok {
		return nil, chatrepo.ErrNoRoomHistory
	}

	return msgs, nil
}
