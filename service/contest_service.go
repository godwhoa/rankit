package service

import (
	"context"
	"fmt"

	"rankit/errors"
	"rankit/postgres"
	"rankit/postgres/sqlgen"
	"rankit/rankit"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/ksuid"
	"go.uber.org/multierr"
)

type ContestService struct {
	querier postgres.Querier
	db      *pgxpool.Pool
}

var _ rankit.ContestService = (*ContestService)(nil)

func NewContestService(querier postgres.Querier) *ContestService {
	return &ContestService{
		querier: querier,
	}
}

func (cs *ContestService) CreateContest(ctx context.Context, p rankit.CreateContestParam) (*rankit.Contest, error) {
	contest, err := cs.querier.CreateContest(ctx, sqlgen.CreateContestParams{
		ID:          generateContestID(),
		CreatorID:   p.CreatorID,
		Title:       p.Title,
		Description: p.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create contest: %w", err)
	}

	return toRankitContest(contest), nil
}

func (cs *ContestService) AddItem(ctx context.Context, p rankit.AddItemParam) (*rankit.Item, error) {
	creatorID, err := cs.querier.GetContestCreatorID(ctx, p.ContestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contest creator id: %w", err)
	}
	if creatorID != p.CreatorID {
		return nil, errors.E(errors.Forbidden, "only the creator of the contest can add items")
	}

	item, err := cs.querier.CreateItem(ctx, sqlgen.CreateItemParams{
		ID:          generateItemID(),
		ContestID:   p.ContestID,
		ContentType: string(p.ContentType),
		Content:     p.Content,
		EloRating:   rankit.DEFAULT_ELO_RATING,
	})
	if _, ok := postgres.IsForeignKeyViolation(err); ok {
		return nil, errors.E(errors.Invalid, "invalid contest id")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	return toRankitItem(item), nil
}

func (cs *ContestService) GetContest(ctx context.Context, id string) (*rankit.Contest, error) {
	contest, err := cs.querier.GetContestByID(ctx, id)
	if postgres.IsNotFound(err) {
		return nil, errors.E(errors.NotFound, "contest not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get contest by id: %w", err)
	}
	rcontest := toRankitContest(contest)

	items, err := cs.querier.GetItemsByContestID(ctx, id)
	if err != nil && postgres.IsNotFound(err) {
		return nil, fmt.Errorf("failed to get items by contest id: %w", err)
	}
	rcontest.Items = Map(items, toRankitItem)
	return rcontest, nil
}

func (cs *ContestService) RecordVote(ctx context.Context, p rankit.RecordVoteParam) error {
	tx, err := cs.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for voting: %w", err)
	}
	defer tx.Rollback(ctx)
	q := cs.querier.WithTx(tx)

	// Fetch the current ratings
	ratings, err := q.GetWinnerLoserEloRatings(ctx, sqlgen.GetWinnerLoserEloRatingsParams{
		WinnerID: p.WinnerItemID,
		LoserID:  p.LoserItemID,
	})
	if postgres.IsNotFound(err) {
		return errors.E(errors.Invalid, "invalid winner or loser item id")
	}
	if err != nil {
		return fmt.Errorf("failed to get winner and loser elo ratings for voting: %w", err)
	}

	// Calculate the new ratings and update the ratings
	var updateErrs error
	winnerRating, loserRating := rankit.CalculateElo(ratings.WinnerEloRating, ratings.LoserEloRating)
	q.UpdateItemsEloRatings(ctx, []sqlgen.UpdateItemsEloRatingsParams{
		{ID: p.WinnerItemID, EloRating: winnerRating},
		{ID: p.LoserItemID, EloRating: loserRating},
	}).Exec(func(i int, err error) {
		updateErrs = multierr.Append(updateErrs, err)
	})
	if updateErrs != nil {
		return fmt.Errorf("failed to update elo ratings for voting: %w", updateErrs)
	}

	// Record Vote
	err = q.RecordVote(ctx, sqlgen.RecordVoteParams{
		ID:           generateVoteID(),
		ContestID:    p.ContestID,
		UserID:       p.VoterID,
		WinnerItemID: p.WinnerItemID,
		LoserItemID:  p.LoserItemID,
	})
	if err != nil {
		return fmt.Errorf("failed to record vote: %w", err)
	}

	// Record elo history
	var recordErrs error
	q.RecordEloHistories(ctx, []sqlgen.RecordEloHistoriesParams{
		{ID: generateEloHistoryID(), ItemID: p.WinnerItemID, EloRating: winnerRating},
		{ID: generateEloHistoryID(), ItemID: p.LoserItemID, EloRating: loserRating},
	}).Exec(func(i int, err error) {
		recordErrs = multierr.Append(recordErrs, err)
	})
	if recordErrs != nil {
		return fmt.Errorf("failed to record elo history for voting: %w", recordErrs)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction for recording vote: %w", err)
	}
	return nil
}

func (cs *ContestService) GetMatchUp(ctx context.Context, contestID string) ([]*rankit.Item, error) {
	items, err := cs.querier.GetRandomItems(ctx, sqlgen.GetRandomItemsParams{
		ContestID: contestID,
		Limit:     2,
	})
	if len(items) != 2 {
		return nil, errors.E(errors.Invalid, "not enough items to create a match up")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get random items: %w", err)
	}

	return Map(items, toRankitItem), nil
}

func toRankitContest(c *sqlgen.Contest) *rankit.Contest {
	return &rankit.Contest{
		ID:          c.ID,
		CreatorID:   c.CreatorID,
		Title:       c.Title,
		Description: c.Description,
		ContentType: rankit.ContentType(c.ContentType),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func toRankitItem(i *sqlgen.Item) *rankit.Item {
	return &rankit.Item{
		ID:          i.ID,
		ContestID:   i.ContestID,
		ContentType: rankit.ContentType(i.ContentType),
		Content:     i.Content,
		EloRating:   i.EloRating,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

func generateContestID() string {
	return "con_" + ksuid.New().String()
}

func generateItemID() string {
	return "itm_" + ksuid.New().String()
}

func generateVoteID() string {
	return "vot_" + ksuid.New().String()
}

func generateEloHistoryID() string {
	return "hst_" + ksuid.New().String()
}

func Map[A any, B any](xs []A, f func(A) B) []B {
	ys := make([]B, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}
