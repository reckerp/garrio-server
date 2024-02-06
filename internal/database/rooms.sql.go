// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: rooms.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createRoom = `-- name: CreateRoom :one
INSERT INTO rooms (name, invite_code, owner_id, record_messages, anon_users) VALUES ($1, $2, $3, $4, $5) RETURNING id, name, invite_code, record_messages, anon_users, owner_id, created_at
`

type CreateRoomParams struct {
	Name           string
	InviteCode     string
	OwnerID        uuid.UUID
	RecordMessages bool
	AnonUsers      bool
}

func (q *Queries) CreateRoom(ctx context.Context, arg CreateRoomParams) (Room, error) {
	row := q.db.QueryRowContext(ctx, createRoom,
		arg.Name,
		arg.InviteCode,
		arg.OwnerID,
		arg.RecordMessages,
		arg.AnonUsers,
	)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.InviteCode,
		&i.RecordMessages,
		&i.AnonUsers,
		&i.OwnerID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteRoomByIDAndOwnerID = `-- name: DeleteRoomByIDAndOwnerID :one
DELETE FROM rooms WHERE id = $1 AND owner_id = $2 RETURNING id, name, invite_code, record_messages, anon_users, owner_id, created_at
`

type DeleteRoomByIDAndOwnerIDParams struct {
	ID      uuid.UUID
	OwnerID uuid.UUID
}

func (q *Queries) DeleteRoomByIDAndOwnerID(ctx context.Context, arg DeleteRoomByIDAndOwnerIDParams) (Room, error) {
	row := q.db.QueryRowContext(ctx, deleteRoomByIDAndOwnerID, arg.ID, arg.OwnerID)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.InviteCode,
		&i.RecordMessages,
		&i.AnonUsers,
		&i.OwnerID,
		&i.CreatedAt,
	)
	return i, err
}

const getRoomByID = `-- name: GetRoomByID :one
SELECT id, name, invite_code, record_messages, anon_users, owner_id, created_at FROM rooms WHERE id = $1
`

func (q *Queries) GetRoomByID(ctx context.Context, id uuid.UUID) (Room, error) {
	row := q.db.QueryRowContext(ctx, getRoomByID, id)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.InviteCode,
		&i.RecordMessages,
		&i.AnonUsers,
		&i.OwnerID,
		&i.CreatedAt,
	)
	return i, err
}

const getRoomByInviteCode = `-- name: GetRoomByInviteCode :one
SELECT id, name, invite_code, record_messages, anon_users, owner_id, created_at FROM rooms WHERE invite_code = $1
`

func (q *Queries) GetRoomByInviteCode(ctx context.Context, inviteCode string) (Room, error) {
	row := q.db.QueryRowContext(ctx, getRoomByInviteCode, inviteCode)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.InviteCode,
		&i.RecordMessages,
		&i.AnonUsers,
		&i.OwnerID,
		&i.CreatedAt,
	)
	return i, err
}

const getRoomsWhereUserIdMember = `-- name: GetRoomsWhereUserIdMember :many
SELECT id, name, invite_code, record_messages, anon_users, owner_id, created_at FROM rooms WHERE id IN (SELECT room_id FROM room_members WHERE user_id = $1)
`

func (q *Queries) GetRoomsWhereUserIdMember(ctx context.Context, userID uuid.UUID) ([]Room, error) {
	rows, err := q.db.QueryContext(ctx, getRoomsWhereUserIdMember, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Room
	for rows.Next() {
		var i Room
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.InviteCode,
			&i.RecordMessages,
			&i.AnonUsers,
			&i.OwnerID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
