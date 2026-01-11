package domain

import "context"

type UserScore struct {
	UserID string
	Score  float64
	Rank   int64
}

type LeaderboardRepository interface {
	AddScore(ctx context.Context, userID string, score float64) error
	GetTopRankers(ctx context.Context, n int64) ([]UserScore, error)
	GetRank(ctx context.Context, userID string) (int64, error)
	BanUser(ctx context.Context, userID string) error
	IsBanned(ctx context.Context, userID string) (bool, error)
}
