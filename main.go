package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// this is a global variable to hold the redis client connection
var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func main() {
	rdb.Del(ctx, "match:queue")
	fmt.Println("Matching start...Finding for your opponent")

	var wg sync.WaitGroup

	for i := 1; i <= 20; i++ {
		wg.Add(1)

		go func(playerID int) {
			defer wg.Done()

			score := rand.Intn(3000)
			err := joinQueue(playerID, score)
			if err != nil {
				fmt.Println(err)
			}
		}(i)
	}
	wg.Wait()

	fmt.Println("\n Current matchmaking queue: (ranked by score)")
	showQueue()
	startMatch()
}

func joinQueue(playerID int, score int) error {
	err := rdb.ZAdd(ctx, "match:queue", redis.Z{
		Score:  float64(score),
		Member: fmt.Sprintf("player:%d", playerID),
	}).Err()

	if err != nil {
		return fmt.Errorf("failed to add player %d: %w", playerID, err)
	}

	fmt.Printf("Player %d with score %d joined the queue\n", playerID, score)
	return nil
}

func showQueue() {
	players := rdb.ZRevRangeWithScores(ctx, "match:queue", 0, -1).Val()

	for _, player := range players {
		fmt.Printf(" %v with score %.0f\n", player.Member, player.Score)
	}
}

func startMatch() {
	fmt.Println("Start matching...")
	roomID := 1
	for {
		players := rdb.ZRange(ctx, "match:queue", 0, 4).Val()
		if len(players) < 5 {
			fmt.Printf("\n Not enough players to start a match... %d players left\n", len(players))
			break
		}
		rdb.ZRem(ctx, "match:queue", players[0], players[1], players[2], players[3], players[4])

		roomKey := fmt.Sprintf("match:room:%d", roomID)
		for _, player := range players {
			rdb.RPush(ctx, roomKey, player)
		}
		fmt.Printf("Room %d created with players: %v\n", roomID, players)
		roomID++
	}
}
