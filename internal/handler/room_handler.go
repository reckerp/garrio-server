package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/reckerp/garrio-server/internal/requestresponse"
	"github.com/reckerp/garrio-server/internal/service"
)

type RoomHandler struct {
	roomService *service.RoomService
}

func NewRoomHandler(roomService *service.RoomService) *RoomHandler {
	return &RoomHandler{roomService: roomService}
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	// Check if the user is authenticated
	ctx, err := autenticateUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Create the room
	var roomRequest requestresponse.RoomCreateRequest
	if err := c.ShouldBindJSON(&roomRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := h.roomService.CreateRoom(ctx, &roomRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating room"})
		return
	}

	c.JSON(http.StatusCreated, room)
}

func (h *RoomHandler) BecomeRoomMember(c *gin.Context) {
	// Check if the user is authenticated
	ctx, err := autenticateUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Become a member of the room
	inviteCode := c.Param("invite_code")
	res, err := h.roomService.BecomeRoomMember(ctx, &inviteCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	// Check if the user is authenticated
	ctx, err := autenticateUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Delete the room
	roomID := uuid.MustParse(c.Param("id"))
	err = h.roomService.DeleteRoom(ctx, &roomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "either this room doesnt exist or you are not the owner"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "room deleted"})
}
