package handler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"os"
	"strings"
)

func autenticateUser(c *gin.Context) (context.Context, error) {
	cookie, err := c.Request.Cookie("garrio_jwt")
	bearer := c.GetHeader("Authorization")
	if err != nil && bearer == "" {
		return nil, errors.New("unauthorized")
	}

	var tokenString string
	if bearer == "" {
		tokenString = cookie.Value
	} else {
		tokenString = strings.Split(bearer, " ")[1]
		tokenString = strings.TrimSpace(tokenString)
	}

	log.Println("Token:" + tokenString)

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
	tokenUName := token.Claims.(jwt.MapClaims)["username"].(string)
	ctx := context.WithValue(c.Request.Context(), "uname", tokenUName)
	ctx = context.WithValue(ctx, "userID", tokenUserID)
	return ctx, nil
}
