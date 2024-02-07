package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/reckerp/garrio-server/internal/database"
	"github.com/reckerp/garrio-server/internal/handler"
	"github.com/reckerp/garrio-server/internal/service"
	"github.com/reckerp/garrio-server/internal/ws"
)

func main() {
	godotenv.Load()

	port := os.Getenv("GARRIO_SERVER_PORT")
	if port == "" {
		log.Fatal("GARRIO_SERVER_PORT environment variable is not set")
	}

	dbURL := os.Getenv("GARRIO_DB_URL")
	if dbURL == "" {
		log.Fatal("GARRIO_DB_URL environment variable is not set")
	}

	fmt.Println("DB URL: ", dbURL)

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	queries := database.New(conn)

	userService := service.NewUserService(queries)
	roomService := service.NewRoomService(queries)

	chatHub := ws.NewHub(roomService)

	userHandler := handler.NewUserHandler(userService)
	roomHandler := handler.NewRoomHandler(roomService)
	chatHandler := handler.NewChatHandler(chatHub, roomService)

	go chatHub.Run()

	mainRouter := gin.Default()

	v1Router := mainRouter.Group("/v1")

	v1Router.GET("/test/:roomId", func(c *gin.Context) {
		roomID := uuid.MustParse(c.Param("roomId"))
		room, err := roomService.GetRoomByID(c.Request.Context(), &roomID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"room": room})
	})

	v1Router.POST("/signup", userHandler.CreateUser)
	v1Router.POST("/login", userHandler.LoginUser)
	v1Router.GET("/logout", userHandler.LogoutUser)

	v1Router.POST("/rooms", roomHandler.CreateRoom)
	v1Router.GET("/rooms/join/:invite_code", roomHandler.BecomeRoomMember)
	v1Router.DELETE("/rooms/:id", roomHandler.DeleteRoom)

	v1Router.GET("/rooms/chat/:room_id", chatHandler.JoinRoom)

	// Start the server
	mainRouter.Run(":" + port)
}
