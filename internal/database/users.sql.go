// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: users.sql

package database

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, username, password, created_at, last_login
`

type CreateUserParams struct {
	Username string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
		&i.LastLogin,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, password, created_at, last_login FROM users WHERE username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
		&i.LastLogin,
	)
	return i, err
}

const updateUserLoginTime = `-- name: UpdateUserLoginTime :exec
UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE username = $1
`

func (q *Queries) UpdateUserLoginTime(ctx context.Context, username string) error {
	_, err := q.db.ExecContext(ctx, updateUserLoginTime, username)
	return err
}