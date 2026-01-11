package redis

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
)

func TestAddScore(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewRedisLeaderboardRepository(db, "leaderboard", "banned")

	ctx := context.Background()
	userID := "user1"
	score := 100.0

	mock.ExpectZAdd("leaderboard", redis.Z{
		Score:  score,
		Member: userID,
	}).SetVal(1)

	err := repo.AddScore(ctx, userID, score)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetTopRankers(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewRedisLeaderboardRepository(db, "leaderboard", "banned")

	ctx := context.Background()
	n := int64(3)

	mock.ExpectZRevRangeWithScores("leaderboard", 0, n-1).SetVal([]redis.Z{
		{Score: 300, Member: "user3"},
		{Score: 200, Member: "user2"},
		{Score: 100, Member: "user1"},
	})

	rankers, err := repo.GetTopRankers(ctx, n)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(rankers) != 3 {
		t.Errorf("expected 3 rankers, got %d", len(rankers))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetRank(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewRedisLeaderboardRepository(db, "leaderboard", "banned")

	ctx := context.Background()
	userID := "user2"

	mock.ExpectZRevRank("leaderboard", userID).SetVal(1)

	rank, err := repo.GetRank(ctx, userID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rank != 2 {
		t.Errorf("expected rank 2, got %d", rank)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestBanUser(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewRedisLeaderboardRepository(db, "leaderboard", "banned")

	ctx := context.Background()
	userID := "cheater1"

	mock.ExpectSAdd("banned", userID).SetVal(1)

	err := repo.BanUser(ctx, userID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestIsBanned(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewRedisLeaderboardRepository(db, "leaderboard", "banned")

	ctx := context.Background()
	userID := "cheater1"

	mock.ExpectSIsMember("banned", userID).SetVal(true)

	banned, err := repo.IsBanned(ctx, userID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !banned {
		t.Errorf("expected user to be banned")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
