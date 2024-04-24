package ws

import (
	"context"

	"github.com/google/uuid"
	"github.com/reckerp/garrio-server/internal/database"
	"github.com/reckerp/garrio-server/internal/service"
)

type Room struct {
	ID             uuid.UUID             `json:"id"`
	Name           string                `json:"name"`
	InviteCode     string                `json:"invite_code"`
	RecordMessages bool                  `json:"record_messages"`
	AllowAnon      bool                  `json:"allow_anon"`
	OwnerID        uuid.UUID             `json:"owner_id"`
	Clients        map[uuid.UUID]*Client `json:"clients"`
}

type Hub struct {
	RoomService  *service.RoomService
	ActiveRooms  map[uuid.UUID]*Room
	RegisterCh   chan *Client
	UnregisterCh chan *Client
	BroadcastCh  chan *Message
}

func NewHub(roomService *service.RoomService) *Hub {
	return &Hub{
		RoomService:  roomService,
		ActiveRooms:  make(map[uuid.UUID]*Room),
		RegisterCh:   make(chan *Client),
		UnregisterCh: make(chan *Client),
		BroadcastCh:  make(chan *Message, 5),
	}
}

func (h *Hub) IsRoomActive(roomId uuid.UUID) *Room {
	return h.ActiveRooms[roomId]
}

func (h *Hub) RegisterRoom(room *Room) {
	h.ActiveRooms[room.ID] = room
}

func (h *Hub) Run() {
	for {
		select {
		case registerClient := <-h.RegisterCh:
			room, ok := h.ActiveRooms[registerClient.RoomID]
			if !ok {
				dbRoom, err := h.RoomService.GetRoomByID(context.Background(), &registerClient.RoomID)
				if err != nil {
					continue
				}
				room = DBRoomToWSRoom(dbRoom)
				h.ActiveRooms[registerClient.RoomID] = room
			}
			if _, ok := room.Clients[registerClient.UserID]; !ok {
				room.Clients[registerClient.UserID] = registerClient
			}

		case unregisterClient := <-h.UnregisterCh:
			room, ok := h.ActiveRooms[unregisterClient.RoomID]
			if ok {
				if _, ok := room.Clients[unregisterClient.UserID]; ok {
					delete(room.Clients, unregisterClient.UserID)
					close(unregisterClient.MessageCh)
					if len(room.Clients) != 0 {
						msg := NewMessage(room.ID, unregisterClient.UserID, unregisterClient.ActiveUname,
							unregisterClient.ActiveUname+" left the room.", DISPLAYABLE_SYSTEM_MESSAGE)
						h.BroadcastCh <- msg
					} else {
						delete(h.ActiveRooms, unregisterClient.RoomID)
					}
				}
			}

		case msg := <-h.BroadcastCh:
			if _, ok := h.ActiveRooms[msg.RoomID]; ok {
				for _, client := range h.ActiveRooms[msg.RoomID].Clients {
					if client.UserID != msg.SenderID || msg.MessageType < 51 {
						client.MessageCh <- msg
					}
				}
			}
		}
	}
}

func DBRoomToWSRoom(dbRoom *database.Room) *Room {
	return &Room{
		ID:             dbRoom.ID,
		Name:           dbRoom.Name,
		InviteCode:     dbRoom.InviteCode,
		RecordMessages: dbRoom.RecordMessages,
		AllowAnon:      dbRoom.AnonUsers,
		OwnerID:        dbRoom.OwnerID,
		Clients:        make(map[uuid.UUID]*Client),
	}
}
