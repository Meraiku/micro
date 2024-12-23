package v1

import (
	"net/http"

	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/websocket/intrenal/services/chat"
)

func (s *ChatServiceAPI) handleGlobalChat(w http.ResponseWriter, r *http.Request) {
	log := logging.L(r.Context())

	username, err := r.Cookie("username")
	if err != nil {
		username.Value = "guest"
	}

	client := chat.NewClient(username.Value)

	log.Info(
		"connecting to global chat",
		logging.String("username", client.Username),
	)

	if err := s.cs.ConnectGlobal(client, w, r); err != nil {
		log.Error(
			"failed to connect global chat",
			logging.Err(err),
			logging.String("username", client.Username),
		)
		return
	}
}
