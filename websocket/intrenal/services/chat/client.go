package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/meraiku/micro/pkg/logging"
	"golang.org/x/sync/errgroup"
)

type Client struct {
	ID       uuid.UUID
	Username string

	ChatRoom *Room

	conn    *websocket.Conn
	recieve chan []byte
}

func NewClient(username string) *Client {
	return &Client{
		ID:       uuid.New(),
		Username: username,
		recieve:  make(chan []byte),
	}
}

func (c *Client) StartSession(ctx context.Context, conn *websocket.Conn) error {

	defer func() {
		c.conn.Close()
	}()

	log := logging.L(ctx)

	log.Info(
		"starting new session",
		logging.String("client_id", c.ID.String()),
		logging.String("username", c.Username),
		logging.String("room_id", c.ChatRoom.ID.String()),
	)

	gr, ctx := errgroup.WithContext(ctx)

	gr.Go(c.read)
	gr.Go(c.write)

	if err := gr.Wait(); err != nil {
		log.Info(
			"session closed",
			logging.String("client_id", c.ID.String()),
		)
		return err
	}

	return nil
}

func (c *Client) addConnection(conn *websocket.Conn) {
	c.conn = conn
}

func (c *Client) read() error {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case msg, ok := <-c.recieve:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return ErrRoomClosed
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return fmt.Errorf("next writer: %v", err)
			}

			w.Write(msg)

			// Add queued chat messages to the current websocket message.
			n := len(c.recieve)
			for i := 0; i < n; i++ {
				w.Write(msg)
			}

			if err := w.Close(); err != nil {
				return fmt.Errorf("close writer: %v", err)
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return ErrClientNotAvailable
			}
		}
	}
}

func (c *Client) write() error {

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, text, err := c.conn.ReadMessage()
		log.Printf("got message: %v", string(text))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("connection closed: %v", err)
			}
			break
		}

		msg := &Message{}

		reader := bytes.NewReader(text)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(msg)
		if err != nil {
			log.Printf("decode msg: %v", err)
		}

		c.ChatRoom.Broadcast <- NewMessage(c, msg.Text)
	}

	return nil
}
