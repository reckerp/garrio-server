package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/reckerp/garrio-server/internal/database"
	"github.com/reckerp/garrio-server/internal/requestresponse"
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

func (s *RoomService) DeleteRoom(c context.Context, roomID *uuid.UUID) error {
	requestUserID := uuid.MustParse(c.Value("userID").(string))
	_, err := s.db.DeleteRoomByIDAndOwnerID(c, database.DeleteRoomByIDAndOwnerIDParams{
		ID:      *roomID,
		OwnerID: requestUserID,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *RoomService) BecomeRoomMember(c context.Context, inviteCode *string) (*requestresponse.RoomMemberResponse, error) {
	userID := uuid.MustParse(c.Value("userID").(string))
	room, err := s.db.GetRoomByInviteCode(c, *inviteCode)
	if err != nil {
		return nil, errors.New("A room with that invite code does not exist")
	}

	_, err = s.db.JoinRoomByRoomIDAndUserID(c, database.JoinRoomByRoomIDAndUserIDParams{
		RoomID: room.ID,
		UserID: userID,
	})

	if err != nil {
		return nil, errors.New("You are already a member of this room")
	}

	res := requestresponse.RoomMemberResponse{
		RoomID:         room.ID,
		RoomName:       room.Name,
		RecordMessages: room.RecordMessages,
		AllowAnon:      room.AnonUsers,
	}

	return &res, nil
}

func (s *RoomService) GetRoomByInviteCode(c context.Context, inviteCode *string) (*database.Room, error) {
	room, err := s.db.GetRoomByInviteCode(c, *inviteCode)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &room, nil
}

func (s *RoomService) GetRoomByID(c context.Context, roomID *uuid.UUID) (*database.Room, error) {
	log.Println("CALLED IN SERVICE GetRoomByID: ", roomID)
	room, err := s.db.GetRoomByID(c, *roomID)
	log.Println("ROOM: ", room)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &room, nil
}

func (s *RoomService) IsUserMemberOfRoom(c context.Context, userID *uuid.UUID, roomID *uuid.UUID) bool {
	_, err := s.db.IsUserMemberOfRoom(c, database.IsUserMemberOfRoomParams{
		RoomID: *roomID,
		UserID: *userID,
	})

	return err == nil

}

func (s *RoomService) RoomMemberCountByRoomID(c context.Context, roomID *uuid.UUID) (int64, error) {
	count, err := s.db.RoomMemberCountByRoomID(c, *roomID)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return count, nil
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
