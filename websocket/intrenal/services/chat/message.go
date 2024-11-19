package chat

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/google/uuid"
	"github.com/meraiku/micro/pkg/logging"
)

type Message struct {
	ID       uuid.UUID
	ClientID uuid.UUID
	Username string
	Text     string
}

func NewMessage(client *Client, text string) *Message {
	msg := &Message{
		ID:       uuid.New(),
		ClientID: client.ID,
		Username: client.Username,
		Text:     text,
	}

	return msg
}

func (m *Message) String() string {
	return fmt.Sprintf("%s: %s", m.Username, m.Text)
}

func (m *Message) Render() []byte {
	tmpl, err := template.ParseFS(templates, "templates/message.html")
	if err != nil {
		logging.Default().Error("failed to parse template", logging.Err(err))
		return nil
	}

	var renderedMsg bytes.Buffer
	err = tmpl.Execute(&renderedMsg, m)
	if err != nil {
		logging.Default().Error("failed to execute template", logging.Err(err))
		return nil
	}

	return renderedMsg.Bytes()

}
