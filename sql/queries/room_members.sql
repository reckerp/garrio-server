-- name: JoinRoomByRoomIDAndUserID :one
INSERT INTO room_members (room_id, user_id, is_admin) VALUES ($1, $2, $3) RETURNING *;

-- name: GetRoomMembersByRoomID :many
SELECT * FROM room_members WHERE room_id = $1;

-- name: RoomMemberCountByRoomID :one
SELECT COUNT(*) FROM room_members WHERE room_id = $1;

-- name: LeaveRoomByRoomIDAndUserID :one
DELETE FROM room_members WHERE room_id = $1 AND user_id = $2 RETURNING *;
