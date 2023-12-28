package rankit

import (
	"context"
	"encoding/json"
	"time"
)

type Contest struct {
	ID          string      `json:"id"`
	CreatorID   string      `json:"creator_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	ContentType ContentType `json:"content_type"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Items       []*Item     `json:"items"`
}

type Item struct {
	ID          string          `json:"id"`
	ContestID   string          `json:"contest_id"`
	ContentType ContentType     `json:"content_type"`
	Content     json.RawMessage `json:"content"`
	EloRating   int             `json:"elo_rating"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type ContentType string

const (
	ContentTypeImage ContentType = "v1/image"
	ContentTypeText  ContentType = "v1/text"
)

type Image struct {
	URL string `json:"url"`
}

type Text struct {
	Text string `json:"text"`
}

type RecordVoteParam struct {
	ContestID    string `json:"contest_id"`
	VoterID      string `json:"voter_id"`
	WinnerItemID string `json:"winner_item_id"`
	LoserItemID  string `json:"loser_item_id"`
}

type CreateContestParam struct {
	CreatorID   string `json:"creator_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type AddItemParam struct {
	ContestID   string          `json:"contest_id"`
	ContentType ContentType     `json:"content_type"`
	Content     json.RawMessage `json:"content"`
	CreatorID   string          `json:"creator_id"`
}

type ContestService interface {
	CreateContest(ctx context.Context, p CreateContestParam) (*Contest, error)
	AddItem(ctx context.Context, p AddItemParam) (*Item, error)
	GetContest(ctx context.Context, id string) (*Contest, error)
	RecordVote(ctx context.Context, p RecordVoteParam) error
	GetMatchUp(ctx context.Context, contestID string) ([]*Item, error)
}
