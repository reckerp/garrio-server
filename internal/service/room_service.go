package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/reckerp/garrio-server/internal/database"
	"github.com/reckerp/garrio-server/internal/requestresponse"
	"log"
	"strings"
)

type RoomService struct {
	db *database.Queries
}

func NewRoomService(db *database.Queries) *RoomService {
	return &RoomService{db: db}
}

func (s *RoomService) CreateRoom(c context.Context, req *requestresponse.RoomCreateRequest) (*database.Room, error) {
	// Create the room
	roomOwner := uuid.MustParse(c.Value("userID").(string))

	room, err := s.db.CreateRoom(c, database.CreateRoomParams{
		Name:           req.Name,
		OwnerID:        roomOwner,
		InviteCode:     generateInviteToken(6),
		RecordMessages: req.RecordMessages,
		AnonUsers:      req.AllowAnon,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &room, nil
}

func (s *RoomService) DeleteRoom(c context.Context, roomID uuid.UUID) error {
	requestUserID := uuid.MustParse(c.Value("userID").(string))
	_, err := s.db.DeleteRoomByIDAndOwnerID(c, database.DeleteRoomByIDAndOwnerIDParams{
		ID:      roomID,
		OwnerID: requestUserID,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func generateInviteToken(length int8) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	token := base64.StdEncoding.EncodeToString(bytes)
	token = strings.TrimRight(token, "=")
	token = strings.ToLower(token)

	return token
}
