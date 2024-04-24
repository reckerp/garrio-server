package ws

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn        *websocket.Conn
	MessageCh   chan *Message
	ActiveUname string    `json:"activeUname"`
	UserID      uuid.UUID `json:"userId"`
	RoomID      uuid.UUID `json:"roomId"`
}

func NewClient(conn *websocket.Conn, roomID uuid.UUID, userID uuid.UUID, activeUname string) *Client {
	return &Client{
		Conn:        conn,
		MessageCh:   make(chan *Message, 10),
		ActiveUname: activeUname,
		RoomID:      roomID,
		UserID:      userID,
	}
}

func (c *Client) SendToClient() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		msg, ok := <-c.MessageCh
		if !ok {
			return
		}
		c.Conn.WriteJSON(msg)
	}

}

func (c *Client) ReadFromClient(hub *Hub) {
	defer func() {
		hub.UnregisterCh <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[ERROR] UnexpectedCloseError: %v", err)
			}
			break
		}

		msg := NewMessage(c.RoomID, c.UserID, c.ActiveUname, string(m), USER_MESSAGE)
		hub.BroadcastCh <- msg
	}
}

const (
	DISPLAYABLE_SYSTEM_MESSAGE = 1
	INTERNAL_SYSTEM_MESSAGE    = 2
	ROOM_UPDATE_MESSAGE        = 3
	USER_JOIN_LEAVE_MESSAGE    = 51
	USER_MESSAGE               = 99
)

type Message struct {
	ID          uuid.UUID `json:"id"`
	RoomID      uuid.UUID `json:"room_id"`
	SenderID    uuid.UUID `json:"sender_id"`
	SenderUname string    `json:"sender_uname"`
	Content     string    `json:"content"`
	MessageType int       `json:"message_type"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewMessage(roomID uuid.UUID, senderID uuid.UUID, senderUname string, content string, messageType int) *Message {
	return &Message{
		ID:          uuid.New(),
		RoomID:      roomID,
		SenderID:    senderID,
		SenderUname: senderUname,
		Content:     content,
		MessageType: messageType,
		CreatedAt:   time.Now(),
	}
}
