package service

import (
	"context"
	"fmt"
	"log"
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
	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		if os.Getenv("GIN_MODE") == "release" {
			log.Fatal(err)
		}
		println(err)
	}
	fmt.Println("redis connection established")
}
func PublishTopic(topic string, payload interface{}) error {
	if err := RedisClient.Publish(Ctx, topic, payload).Err(); err != nil {
		return err
		// return err
	}
	return nil
}
