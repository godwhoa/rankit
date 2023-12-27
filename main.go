package main

// import (
// 	"context"
// 	"rankit/postgres/sqlgen"
// 	"strings"

// 	"github.com/jackc/pgx/v5"
// )

// const DEFAULT_ELO_RATING = 1200

// func App() error {
// 	ctx := context.Background()
// 	db, err := pgx.Connect(ctx, "postgres://rankit:rankit@localhost:5432/rankit")
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close(ctx)

// 	tx, err := db.Begin(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback(ctx)

// 	q := sqlgen.New(db).WithTx(tx)

// 	// Fetch the current ratings
// 	ratings, _ := q.GetWinnerLoserEloRatings(ctx, sqlgen.GetWinnerLoserEloRatingsParams{
// 		WinnerID: "1",
// 		LoserID:  "2",
// 	})

// 	// Calculate the new ratings
// 	winnerRating, loserRating := CalculateElo(ratings.WinnerEloRating, ratings.LoserEloRating)

// 	// Update the ratings
// 	var ratingsErr MultiErr
// 	q.UpdateItemsEloRatings(ctx, []sqlgen.UpdateItemsEloRatingsParams{
// 		{ID: "1", EloRating: winnerRating},
// 		{ID: "2", EloRating: loserRating},
// 	}).Exec(collectErrors(ratingsErr))

// 	// Insert Vote record
// 	q.RecordVote(ctx, sqlgen.RecordVoteParams{
// 		ID:           "1",
// 		ContestID:    "1",
// 		UserID:       "1",
// 		WinnerItemID: "1",
// 		LoserItemID:  "2",
// 	})

// 	// Log elo history
// 	var historiesErr MultiErr
// 	q.RecordEloHistories(ctx, []sqlgen.RecordEloHistoriesParams{
// 		{ID: "1", ItemID: "1", EloRating: winnerRating},
// 		{ID: "2", ItemID: "2", EloRating: loserRating},
// 	}).Exec(collectErrors(historiesErr))

// 	tx.Commit(ctx)

// 	return nil
// }

// func collectErrors(errs MultiErr) func(int, error) {
// 	return func(i int, err error) {
// 		if err != nil {
// 			errs = append(errs, err)
// 		}
// 	}
// }

// type MultiErr []error

// func (me MultiErr) Error() string {
// 	msg := &strings.Builder{}
// 	for i, err := range me {
// 		if i > 0 {
// 			msg.WriteString("; ")
// 		}
// 		msg.WriteString(err.Error())
// 	}
// 	return msg.String()
// }

// func main() {
// 	if err := App(); err != nil {
// 		panic(err)
// 	}
// }
