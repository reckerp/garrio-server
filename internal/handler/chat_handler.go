package handler

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	// "github.com/reckerp/garrio-server/internal/requestresponse"
	"github.com/reckerp/garrio-server/internal/requestresponse"
	"github.com/reckerp/garrio-server/internal/service"
	"github.com/reckerp/garrio-server/internal/ws"
)

type ChatHandler struct {
	hub         *ws.Hub
	roomService *service.RoomService
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		//TOOD: Add origin check
		return true
	},
}

func NewChatHandler(hub *ws.Hub, roomService *service.RoomService) *ChatHandler {
	return &ChatHandler{hub: hub, roomService: roomService}
}

func (h *ChatHandler) JoinRoom(c *gin.Context) {
	roomID := uuid.MustParse(c.Param("room_id"))
	ctx, authErr := autenticateUser(c)
	if authErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	wantJoinAsAnon, convErr := strconv.ParseBool(c.GetHeader("JoinAsAnon"))
	if convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JoinAsAnon header"})
		return
	}

	room := h.hub.IsRoomActive(roomID)

	if room == nil {
		dbRoom, dbErr := h.roomService.GetRoomByID(c.Request.Context(), &roomID)
		if dbErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "room not found"})
			return
		}
		room = ws.DBRoomToWSRoom(dbRoom)
		h.hub.RegisterRoom(room)
	}

	if !room.AllowAnon && wantJoinAsAnon {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized - room does not allow anonymous users"})
		return
	}

	userId := uuid.MustParse(ctx.Value("userID").(string))

	isUserMember := h.roomService.IsUserMemberOfRoom(ctx, &userId, &roomID)
	if !isUserMember {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized - not a member of the room"})
		return
	}

	// upgrade the user connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not upgrade connection"})
		return
	}

	clientUname := ctx.Value("uname").(string)
	if wantJoinAsAnon {
		clientUname = genRandomUsername()
	}

	client := ws.NewClient(conn, room.ID, userId, clientUname)

	// Send room update message to client
	roomMemberCount, memberCountErr := h.roomService.RoomMemberCountByRoomID(ctx, &roomID)
	if memberCountErr != nil {
		roomMemberCount = -1
	}
	activeRoomMemberCount := len(room.Clients)
	roomUpdateResponse := requestresponse.RoomUpdateResponse{
		RoomID:          room.ID,
		Name:            room.Name,
		RecordMessages:  room.RecordMessages,
		AllowAnon:       room.AllowAnon,
		ActiveUserCount: int64(activeRoomMemberCount + 1), // +1 for the new user
		MemberCount:     roomMemberCount,
	}
	jsonRoomUpdateResponse, _ := json.Marshal(roomUpdateResponse)
	roomUpdateMsg := ws.NewMessage(room.ID, userId, clientUname, string(jsonRoomUpdateResponse), ws.ROOM_UPDATE_MESSAGE)
	client.MessageCh <- roomUpdateMsg

	// Create a message for the room that the user joined
	msg := ws.NewMessage(room.ID, userId, clientUname, clientUname+" joined the room.", ws.USER_JOIN_LEAVE_MESSAGE)

	// Register client through register channel
	h.hub.RegisterCh <- client
	// Broadcast message to room that user joined
	h.hub.BroadcastCh <- msg

	// Start the client read and write loops
	go client.SendToClient()
	go client.ReadFromClient(h.hub)

}

func genRandomUsername() string {
	bytes := make([]byte, 6)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	identifier := base64.StdEncoding.EncodeToString(bytes)
	identifier = strings.TrimRight(identifier, "=")
	identifier = strings.ToLower(identifier)

	return "AnonymousUser-" + identifier

}
