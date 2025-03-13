// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package sqlc

import (
	"context"
)

const readUserByHandle = `-- name: ReadUserByHandle :one

SELECT id, magic_link, is_active, is_admin
FROM users
WHERE handle = ?1
`

type ReadUserByHandleRow struct {
	ID        int64
	MagicLink string
	IsActive  int64
	IsAdmin   int64
}

//	Copyright (c) 2025 Michael D Henderson. All rights reserved.
//
// ReadUserByHandle gets a user by its handle.
func (q *Queries) ReadUserByHandle(ctx context.Context, handle string) (ReadUserByHandleRow, error) {
	row := q.db.QueryRowContext(ctx, readUserByHandle, handle)
	var i ReadUserByHandleRow
	err := row.Scan(
		&i.ID,
		&i.MagicLink,
		&i.IsActive,
		&i.IsAdmin,
	)
	return i, err
}

const readUserByID = `-- name: ReadUserByID :one
SELECT id, handle, is_active, is_admin
FROM users
WHERE id = ?1
`

type ReadUserByIDRow struct {
	ID       int64
	Handle   string
	IsActive int64
	IsAdmin  int64
}

// ReadUserByID gets a user by its id.
func (q *Queries) ReadUserByID(ctx context.Context, id int64) (ReadUserByIDRow, error) {
	row := q.db.QueryRowContext(ctx, readUserByID, id)
	var i ReadUserByIDRow
	err := row.Scan(
		&i.ID,
		&i.Handle,
		&i.IsActive,
		&i.IsAdmin,
	)
	return i, err
}

const readUserByMagicKey = `-- name: ReadUserByMagicKey :one
SELECT id, is_active, is_admin
FROM users
WHERE handle = ?1
  AND magic_link = ?2
`

type ReadUserByMagicKeyParams struct {
	Handle   string
	MagicKey string
}

type ReadUserByMagicKeyRow struct {
	ID       int64
	IsActive int64
	IsAdmin  int64
}

// ReadUserByMagicKey gets a user by its magic key.
func (q *Queries) ReadUserByMagicKey(ctx context.Context, arg ReadUserByMagicKeyParams) (ReadUserByMagicKeyRow, error) {
	row := q.db.QueryRowContext(ctx, readUserByMagicKey, arg.Handle, arg.MagicKey)
	var i ReadUserByMagicKeyRow
	err := row.Scan(&i.ID, &i.IsActive, &i.IsAdmin)
	return i, err
}
