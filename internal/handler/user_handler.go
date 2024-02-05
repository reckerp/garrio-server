package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/reckerp/garrio-server/internal/requestresponse"
	"github.com/reckerp/garrio-server/internal/service"
	"log"
	"net/http"
	"os"
	"strconv"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var userRequest requestresponse.UserCreateRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	usr, err := h.userService.CreateUser(c.Request.Context(), &userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, usr)
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var userRequest requestresponse.UserLoginRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	usr, err := h.userService.LoginUser(c.Request.Context(), &userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	useHttpsValue := os.Getenv("GARRIO_USE_HTTPS")
	log.Println("GARRIO_USE_HTTPS: ", useHttpsValue)
	useHttps, err := strconv.ParseBool(useHttpsValue)
	if err != nil {
		useHttps = false
		log.Fatal("Error parsing GARRIO_USE_HTTPS")
	}
	c.SetCookie("garrio_jwt", usr.AccessToken, 3600, "/", os.Getenv("GARRIO_HOST"), useHttps, !useHttps)

	c.JSON(http.StatusOK, requestresponse.UserResponse{
		Message:  "user successfully logged in",
		ID:       usr.ID,
		Username: usr.Username,
	})
}

func (h *UserHandler) LogoutUser(c *gin.Context) {
	c.SetCookie("garrio_jwt", "", -1, "/", os.Getenv("GARRIO_HOST"), false, false)
	c.JSON(http.StatusOK, gin.H{"message": "user successfully logged out"})
}
