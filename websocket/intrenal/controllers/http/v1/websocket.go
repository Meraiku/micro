package v1

import "net/http"

func (s *ChatServiceAPI) handleWebsocket(w http.ResponseWriter, r *http.Request) {

	serveWs(s.hub, w, r)
}
