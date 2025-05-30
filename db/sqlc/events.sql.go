// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: events.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createFavourite = `-- name: CreateFavourite :one
INSERT INTO "favourites" (
  event_id,
  added_by
) VALUES (
  $1, $2
)
RETURNING id, event_id, added_by, created_at
`

type CreateFavouriteParams struct {
	EventID uuid.NullUUID `json:"event_id"`
	AddedBy uuid.NullUUID `json:"added_by"`
}

func (q *Queries) CreateFavourite(ctx context.Context, arg CreateFavouriteParams) (Favourites, error) {
	row := q.db.QueryRowContext(ctx, createFavourite, arg.EventID, arg.AddedBy)
	var i Favourites
	err := row.Scan(
		&i.ID,
		&i.EventID,
		&i.AddedBy,
		&i.CreatedAt,
	)
	return i, err
}

const deleteFavourite = `-- name: DeleteFavourite :exec
DELETE FROM favourites
WHERE id = $1
`

func (q *Queries) DeleteFavourite(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteFavourite, id)
	return err
}

const deleteFavouriteByUserAndEvent = `-- name: DeleteFavouriteByUserAndEvent :exec
DELETE FROM favourites
WHERE event_id = $1 AND added_by = $2
`

type DeleteFavouriteByUserAndEventParams struct {
	EventID uuid.NullUUID `json:"event_id"`
	AddedBy uuid.NullUUID `json:"added_by"`
}

func (q *Queries) DeleteFavouriteByUserAndEvent(ctx context.Context, arg DeleteFavouriteByUserAndEventParams) error {
	_, err := q.db.ExecContext(ctx, deleteFavouriteByUserAndEvent, arg.EventID, arg.AddedBy)
	return err
}

const getFavourite = `-- name: GetFavourite :one
SELECT id, event_id, added_by, created_at FROM favourites
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetFavourite(ctx context.Context, id uuid.UUID) (Favourites, error) {
	row := q.db.QueryRowContext(ctx, getFavourite, id)
	var i Favourites
	err := row.Scan(
		&i.ID,
		&i.EventID,
		&i.AddedBy,
		&i.CreatedAt,
	)
	return i, err
}

const listFavouritesByEvent = `-- name: ListFavouritesByEvent :many
SELECT id, event_id, added_by, created_at FROM favourites
WHERE event_id = $1
ORDER BY created_at DESC
`

func (q *Queries) ListFavouritesByEvent(ctx context.Context, eventID uuid.NullUUID) ([]Favourites, error) {
	rows, err := q.db.QueryContext(ctx, listFavouritesByEvent, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Favourites
	for rows.Next() {
		var i Favourites
		if err := rows.Scan(
			&i.ID,
			&i.EventID,
			&i.AddedBy,
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

const listFavouritesByUser = `-- name: ListFavouritesByUser :many
SELECT id, event_id, added_by, created_at FROM favourites
WHERE added_by = $1
ORDER BY created_at DESC
`

func (q *Queries) ListFavouritesByUser(ctx context.Context, addedBy uuid.NullUUID) ([]Favourites, error) {
	rows, err := q.db.QueryContext(ctx, listFavouritesByUser, addedBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Favourites
	for rows.Next() {
		var i Favourites
		if err := rows.Scan(
			&i.ID,
			&i.EventID,
			&i.AddedBy,
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
