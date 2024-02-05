package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/reckerp/garrio-server/internal/database"
	"github.com/reckerp/garrio-server/internal/handler"
	"github.com/reckerp/garrio-server/internal/service"
	"log"
	"os"
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
	userHandler := handler.NewUserHandler(userService)

	mainRouter := gin.Default()

	v1Router := mainRouter.Group("/v1")

	v1Router.POST("/signup", userHandler.CreateUser)
	v1Router.POST("/login", userHandler.LoginUser)
	v1Router.GET("/logout", userHandler.LogoutUser)

	// Start the server
	mainRouter.Run(":" + port)
}
