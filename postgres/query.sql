-- Find a user by ID
-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- Find a user by email
-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- Create a new user
-- name: CreateUser :one
INSERT INTO users (
    id, display_name, email, password_hash, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, CURRENT_TIMESTAMP AT TIME ZONE 'UTC', CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
) RETURNING *;

-- Create a new contest
-- name: CreateContest :one
INSERT INTO contests (
    id, creator_id, title, description, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, CURRENT_TIMESTAMP AT TIME ZONE 'UTC', CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
) RETURNING *;

-- Find contest by ID
-- name: GetContestByID :one
SELECT * FROM contests WHERE id = $1;

-- Get creator id of a contest
-- name: GetContestCreatorID :one
SELECT creator_id FROM contests WHERE id = $1;

-- Create a new item
-- name: CreateItem :one
INSERT INTO items (
    id, contest_id, content_type, content, elo_rating, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, CURRENT_TIMESTAMP AT TIME ZONE 'UTC', CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
) RETURNING *;

-- Find items by contest ID
-- name: GetItemsByContestID :many
SELECT * FROM items WHERE contest_id = $1;

-- Get an item's current ELO rating
-- name: GetItemEloRating :one
SELECT elo_rating FROM items WHERE id = $1;

-- Get ELO rating of winner and loser items as a tuple
-- name: GetWinnerLoserEloRatings :one
SELECT
    winner.elo_rating as winner_elo_rating,
    loser.elo_rating as loser_elo_rating
FROM
    items AS winner,
    items AS loser
WHERE
    winner.id = sqlc.arg(winner_id)
    AND loser.id = sqlc.arg(loser_id);

-- Get random items for a contest
-- name: GetRandomItems :many
SELECT * FROM items WHERE contest_id = $1 ORDER BY RANDOM() LIMIT $2;


-- Update an item's ELO rating
-- name: UpdateItemEloRating :exec
UPDATE items SET elo_rating = $2, updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC' WHERE id = $1;

-- Update multiple items' ELO ratings
-- name: UpdateItemsEloRatings :batchexec
UPDATE items SET elo_rating = $2, updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC' WHERE id = $1;

-- Record a vote
-- name: RecordVote :exec
INSERT INTO votes (
    id, contest_id, user_id, winner_item_id, loser_item_id, created_at
) VALUES (
    $1, $2, $3, $4, $5, CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
);

-- Record an item's ELO rating history
-- name: RecordEloHistory :exec
INSERT INTO elo_history (
    id, item_id, elo_rating, created_at
) VALUES (
    $1, $2, $3, CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
);

-- Record multiple items' ELO rating histories
-- name: RecordEloHistories :batchexec
INSERT INTO elo_history (
    id, item_id, elo_rating, created_at
) VALUES (
    $1, $2, $3, CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
);

-- Get item ELO rating history
-- name: GetItemEloHistory :many
SELECT * FROM elo_history WHERE item_id = $1 ORDER BY created_at DESC;