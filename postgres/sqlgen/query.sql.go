// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: query.sql

package sqlgen

import (
	"context"
	"encoding/json"
)

const createContest = `-- name: CreateContest :one
INSERT INTO contests (
    id, creator_id, title, description, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, CURRENT_TIMESTAMP AT TIME ZONE 'UTC', CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
) RETURNING id, creator_id, title, description, content_type, created_at, updated_at
`

type CreateContestParams struct {
	ID          string `json:"id"`
	CreatorID   string `json:"creator_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Create a new contest
func (q *Queries) CreateContest(ctx context.Context, arg CreateContestParams) (*Contest, error) {
	row := q.db.QueryRow(ctx, createContest,
		arg.ID,
		arg.CreatorID,
		arg.Title,
		arg.Description,
	)
	var i Contest
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.Title,
		&i.Description,
		&i.ContentType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const createItem = `-- name: CreateItem :one
INSERT INTO items (
    id, contest_id, content_type, content, elo_rating, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, CURRENT_TIMESTAMP AT TIME ZONE 'UTC', CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
) RETURNING id, contest_id, content_type, content, elo_rating, created_at, updated_at
`

type CreateItemParams struct {
	ID          string          `json:"id"`
	ContestID   string          `json:"contest_id"`
	ContentType string          `json:"content_type"`
	Content     json.RawMessage `json:"content"`
	EloRating   int             `json:"elo_rating"`
}

// Create a new item
func (q *Queries) CreateItem(ctx context.Context, arg CreateItemParams) (*Item, error) {
	row := q.db.QueryRow(ctx, createItem,
		arg.ID,
		arg.ContestID,
		arg.ContentType,
		arg.Content,
		arg.EloRating,
	)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.ContestID,
		&i.ContentType,
		&i.Content,
		&i.EloRating,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    id, display_name, email, password_hash, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, CURRENT_TIMESTAMP AT TIME ZONE 'UTC', CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
) RETURNING id, display_name, email, password_hash, created_at, updated_at
`

type CreateUserParams struct {
	ID           string `json:"id"`
	DisplayName  string `json:"display_name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

// Create a new user
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (*User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.ID,
		arg.DisplayName,
		arg.Email,
		arg.PasswordHash,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.DisplayName,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getContestByID = `-- name: GetContestByID :one
SELECT id, creator_id, title, description, content_type, created_at, updated_at FROM contests WHERE id = $1
`

// Find contest by ID
func (q *Queries) GetContestByID(ctx context.Context, id string) (*Contest, error) {
	row := q.db.QueryRow(ctx, getContestByID, id)
	var i Contest
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.Title,
		&i.Description,
		&i.ContentType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getItemEloHistory = `-- name: GetItemEloHistory :many
SELECT id, item_id, elo_rating, created_at FROM elo_history WHERE item_id = $1 ORDER BY created_at DESC
`

// Get item ELO rating history
func (q *Queries) GetItemEloHistory(ctx context.Context, itemID string) ([]*EloHistory, error) {
	rows, err := q.db.Query(ctx, getItemEloHistory, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*EloHistory
	for rows.Next() {
		var i EloHistory
		if err := rows.Scan(
			&i.ID,
			&i.ItemID,
			&i.EloRating,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getItemEloRating = `-- name: GetItemEloRating :one
SELECT elo_rating FROM items WHERE id = $1
`

// Get an item's current ELO rating
func (q *Queries) GetItemEloRating(ctx context.Context, id string) (int, error) {
	row := q.db.QueryRow(ctx, getItemEloRating, id)
	var elo_rating int
	err := row.Scan(&elo_rating)
	return elo_rating, err
}

const getItemsByContestID = `-- name: GetItemsByContestID :many
SELECT id, contest_id, content_type, content, elo_rating, created_at, updated_at FROM items WHERE contest_id = $1
`

// Find items by contest ID
func (q *Queries) GetItemsByContestID(ctx context.Context, contestID string) ([]*Item, error) {
	rows, err := q.db.Query(ctx, getItemsByContestID, contestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(
			&i.ID,
			&i.ContestID,
			&i.ContentType,
			&i.Content,
			&i.EloRating,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRandomItems = `-- name: GetRandomItems :many
SELECT id, contest_id, content_type, content, elo_rating, created_at, updated_at FROM items WHERE contest_id = $1 ORDER BY RANDOM() LIMIT $2
`

type GetRandomItemsParams struct {
	ContestID string `json:"contest_id"`
	Limit     int32  `json:"limit"`
}

// Get random items for a contest
func (q *Queries) GetRandomItems(ctx context.Context, arg GetRandomItemsParams) ([]*Item, error) {
	rows, err := q.db.Query(ctx, getRandomItems, arg.ContestID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(
			&i.ID,
			&i.ContestID,
			&i.ContentType,
			&i.Content,
			&i.EloRating,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, display_name, email, password_hash, created_at, updated_at FROM users WHERE email = $1
`

// Find a user by email
func (q *Queries) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.DisplayName,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, display_name, email, password_hash, created_at, updated_at FROM users WHERE id = $1
`

// Find a user by ID
func (q *Queries) GetUserByID(ctx context.Context, id string) (*User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.DisplayName,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getWinnerLoserEloRatings = `-- name: GetWinnerLoserEloRatings :one
SELECT
    winner.elo_rating as winner_elo_rating,
    loser.elo_rating as loser_elo_rating
FROM
    items AS winner,
    items AS loser
WHERE
    winner.id = $1
    AND loser.id = $2
`

type GetWinnerLoserEloRatingsParams struct {
	WinnerID string `json:"winner_id"`
	LoserID  string `json:"loser_id"`
}

type GetWinnerLoserEloRatingsRow struct {
	WinnerEloRating int `json:"winner_elo_rating"`
	LoserEloRating  int `json:"loser_elo_rating"`
}

// Get ELO rating of winner and loser items as a tuple
func (q *Queries) GetWinnerLoserEloRatings(ctx context.Context, arg GetWinnerLoserEloRatingsParams) (*GetWinnerLoserEloRatingsRow, error) {
	row := q.db.QueryRow(ctx, getWinnerLoserEloRatings, arg.WinnerID, arg.LoserID)
	var i GetWinnerLoserEloRatingsRow
	err := row.Scan(&i.WinnerEloRating, &i.LoserEloRating)
	return &i, err
}

const recordEloHistory = `-- name: RecordEloHistory :exec
INSERT INTO elo_history (
    id, item_id, elo_rating, created_at
) VALUES (
    $1, $2, $3, CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
)
`

type RecordEloHistoryParams struct {
	ID        string `json:"id"`
	ItemID    string `json:"item_id"`
	EloRating int    `json:"elo_rating"`
}

// Record an item's ELO rating history
func (q *Queries) RecordEloHistory(ctx context.Context, arg RecordEloHistoryParams) error {
	_, err := q.db.Exec(ctx, recordEloHistory, arg.ID, arg.ItemID, arg.EloRating)
	return err
}

const recordVote = `-- name: RecordVote :exec
INSERT INTO votes (
    id, contest_id, user_id, winner_item_id, loser_item_id, created_at
) VALUES (
    $1, $2, $3, $4, $5, CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
)
`

type RecordVoteParams struct {
	ID           string `json:"id"`
	ContestID    string `json:"contest_id"`
	UserID       string `json:"user_id"`
	WinnerItemID string `json:"winner_item_id"`
	LoserItemID  string `json:"loser_item_id"`
}

// Record a vote
func (q *Queries) RecordVote(ctx context.Context, arg RecordVoteParams) error {
	_, err := q.db.Exec(ctx, recordVote,
		arg.ID,
		arg.ContestID,
		arg.UserID,
		arg.WinnerItemID,
		arg.LoserItemID,
	)
	return err
}

const updateItemEloRating = `-- name: UpdateItemEloRating :exec
UPDATE items SET elo_rating = $2, updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC' WHERE id = $1
`

type UpdateItemEloRatingParams struct {
	ID        string `json:"id"`
	EloRating int    `json:"elo_rating"`
}

// Update an item's ELO rating
func (q *Queries) UpdateItemEloRating(ctx context.Context, arg UpdateItemEloRatingParams) error {
	_, err := q.db.Exec(ctx, updateItemEloRating, arg.ID, arg.EloRating)
	return err
}
