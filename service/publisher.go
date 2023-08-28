package service

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var (
	RedisClient *redis.Client
)

func InitRedisClient() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println("redis connection established")
}
func PublishTopic(topic string, payload interface{}) error {
	if err := RedisClient.Publish(Ctx, topic, payload).Err(); err != nil {
		return err
		// return err
	}
	return nil
}
