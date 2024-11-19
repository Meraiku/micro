package chat

import "github.com/google/uuid"

type Message struct {
	ID       uuid.UUID
	ClientID uuid.UUID
	Username string
	Text     string
}
