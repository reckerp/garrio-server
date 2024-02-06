-- name: CreateRoom :one
INSERT INTO rooms (name, invite_code, owner_id, record_messages, anon_users) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetRoomByID :one
SELECT * FROM rooms WHERE id = $1;

-- name: GetRoomByInviteCode :one
SELECT * FROM rooms WHERE invite_code = $1;

-- name: GetRoomsWhereUserIdMember :many
SELECT * FROM rooms WHERE id IN (SELECT room_id FROM room_members WHERE user_id = $1);

-- name: DeleteRoomByIDAndOwnerID :one
DELETE FROM rooms WHERE id = $1 AND owner_id = $2 RETURNING *;
