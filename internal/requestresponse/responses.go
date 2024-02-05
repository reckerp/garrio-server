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
