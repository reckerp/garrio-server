package requestresponse

import (
	"github.com/google/uuid"
	"github.com/reckerp/garrio-server/internal/database"
	"time"
)

type UserCreatedResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUserCreatedResponseFromUser(user *database.User) *UserCreatedResponse {
	res := UserCreatedResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}

	return &res
}

type UserLoginResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	AccessToken string    `json:"access_token"`
}

func NewUserLoginResponseFromUser(user *database.User, accessToken string) *UserLoginResponse {
	res := UserLoginResponse{
		ID:          user.ID,
		Username:    user.Username,
		AccessToken: accessToken,
	}

	return &res
}

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Message  string    `json:"message"`
}

type JoinRoomResponse struct {
	RoomID         uuid.UUID `json:"room_id"`
	RecordMessages bool      `json:"record_messages"`
	ActiveUsers    int       `json:"active_users"`
}

type RoomMemberResponse struct {
	RoomID         uuid.UUID `json:"room_id"`
	RoomName       string    `json:"room_name"`
	RecordMessages bool      `json:"record_messages"`
	AllowAnon      bool      `json:"allow_anon"`
}

type RoomUpdateResponse struct {
	RoomID          uuid.UUID `json:"room_id"`
	Name            string    `json:"name"`
	RecordMessages  bool      `json:"record_messages"`
	AllowAnon       bool      `json:"allow_anon"`
	ActiveUserCount int64     `json:"active_users"`
	MemberCount     int64     `json:"member_count"`
}
