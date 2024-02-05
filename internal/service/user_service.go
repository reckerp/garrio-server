package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/reckerp/garrio-server/internal/database"
	"github.com/reckerp/garrio-server/internal/requestresponse"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

type UserService struct {
	db *database.Queries
}

func NewUserService(db *database.Queries) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(c *gin.Context, user *requestresponse.UserCreateRequest) (*requestresponse.UserCreatedResponse, error) {
	if user.Username == "" || user.Password == "" {
		return nil, errors.New("username and password are required")
	}

	if utf8.RuneCountInString(user.Username) < 4 {
		return nil, errors.New("username must be at least 4 characters")
	}
	if utf8.RuneCountInString(user.Password) < 16 {
		return nil, errors.New("password must be at least 16 characters")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return nil, errors.New("error hashing password")
	}

	usr, err := s.db.CreateUser(c, database.CreateUserParams{
		Username: user.Username,
		Password: hashedPassword,
	})
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return nil, errors.New("username already exists")
	} else if err != nil {
		return nil, err
	}
	user_response := requestresponse.NewUserCreatedResponseFromUser(&usr)
	return user_response, nil
}

func (s *UserService) LoginUser(c *gin.Context, user *requestresponse.UserLoginRequest) (*requestresponse.UserLoginResponse, error) {
	usr, err := s.db.GetUserByUsername(c, user.Username)
	if err != nil {
		return nil, errors.New("username or password incorrect")
	}

	if !checkPasswordHash(user.Password, usr.Password) {
		return nil, errors.New("username or password incorrect")
	}

	accessToken, err := generateAccessToken(&usr)
	if err != nil {
		return nil, errors.New("error generating access token")
	}

	user_response := requestresponse.NewUserLoginResponseFromUser(&usr, accessToken)

	s.db.UpdateUserLoginTime(c, user.Username)

	return user_response, nil
}

func generateAccessToken(user *database.User) (string, error) {
	claims := struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		jwt.RegisteredClaims
	}{
		ID:       user.ID.String(),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "garrio",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("GARRIO_JWT_SECRET")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return signedToken, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
