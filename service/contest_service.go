package service

import (
	"context"
	"rankit/postgres/sqlgen"
	"rankit/rankit"
)

type ContestService struct {
	querier sqlgen.Querier
}

var _ rankit.ContestService = (*ContestService)(nil)

func NewContestService(querier sqlgen.Querier) *ContestService {
	return &ContestService{
		querier: querier,
	}
}

func (cs *ContestService) CreateContest(ctx context.Context, p rankit.CreateContestParam) (*rankit.Contest, error) {
	// TODO: Implement CreateContest method
	return nil, nil
}

func (cs *ContestService) AddItem(ctx context.Context, p rankit.AddItemParam) (*rankit.Item, error) {
	// TODO: Implement AddItem method
	return nil, nil
}

func (cs *ContestService) GetContest(ctx context.Context, id string) (*rankit.Contest, error) {
	// TODO: Implement GetContest method
	return nil, nil
}

func (cs *ContestService) RecordVote(ctx context.Context, p rankit.RecordVoteParam) error {
	// TODO: Implement RecordVote method
	return nil
}

func (cs *ContestService) GetMatchUp(ctx context.Context, contestID string) ([]*rankit.Item, error) {
	// TODO: Implement GetMatchUp method
	return nil, nil
}
