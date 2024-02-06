package handler

import (
	"context"
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func autenticateUser(c *gin.Context) (context.Context, error) {
	cookie, err := c.Request.Cookie("garrio_jwt")
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	tokenString := cookie.Value
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("GARRIO_JWT_SECRET")), nil
	})

	if err != nil {
		return nil, errors.New("unauthorized")
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("unauthorized")
	}

	tokenUserID := token.Claims.(jwt.MapClaims)["id"].(string)
	return context.WithValue(c.Request.Context(), "userID", tokenUserID), nil
}
