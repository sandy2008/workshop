package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	infra "github.com/sokoide/workshop/infra/assets/redis_leaderboard/infra/redis"
	"github.com/sokoide/workshop/infra/assets/redis_leaderboard/usecase"
)

func main() {
	// 1. Setup Redis Client
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// 2. Dependency Injection (DI)
	// Framework/Main layer connects everything
	repo := infra.NewRedisLeaderboardRepository(client, "game_leaderboard", "banned_users")
	uc := usecase.NewLeaderboardUsecase(repo)

	// 3. Command Parsing
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	ctx := context.Background()
	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) != 4 {
			fmt.Println("Usage: add <user_id> <score>")
			return
		}
		userID := os.Args[2]
		score, err := strconv.ParseFloat(os.Args[3], 64)
		if err != nil {
			log.Fatal("Invalid score")
		}
		if err := uc.AddScore(ctx, userID, score); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Added score %.2f for user %s\n", score, userID)

	case "top":
		n := int64(10)
		if len(os.Args) == 3 {
			parsedN, err := strconv.ParseInt(os.Args[2], 10, 64)
			if err == nil {
				n = parsedN
			}
		}
		rankers, err := uc.GetTopRankers(ctx, n)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("--- Top %d Rankers ---\n", n)
		for _, r := range rankers {
			fmt.Printf("%d. %s: %.2f\n", r.Rank, r.UserID, r.Score)
		}

	case "rank":
		if len(os.Args) != 3 {
			fmt.Println("Usage: rank <user_id>")
			return
		}
		userID := os.Args[2]
		rank, err := uc.GetRank(ctx, userID)
		if err != nil {
			log.Fatal(err)
		}
		if rank == 0 {
			fmt.Printf("User %s not found in leaderboard\n", userID)
		} else {
			fmt.Printf("User %s is ranked #%d\n", userID, rank)
		}

	case "ban":
		if len(os.Args) != 3 {
			fmt.Println("Usage: ban <user_id>")
			return
		}
		userID := os.Args[2]
		if err := uc.BanUser(ctx, userID); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("User %s has been banned\n", userID)

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Leaderboard CLI usage:")
	fmt.Println("  add <user_id> <score>  - Add or update user score")
	fmt.Println("  top [n]               - Show top n rankers (default 10)")
	fmt.Println("  rank <user_id>        - Show rank of a specific user")
	fmt.Println("  ban <user_id>         - Ban a user from the leaderboard")
}
