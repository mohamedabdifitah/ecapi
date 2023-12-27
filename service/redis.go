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
		Password: os.Getenv("REDIS_PASS"), // no password set
		DB:       0,                       // use default DB
	})
	if os.Getenv("APP_ENV") != "development" {
		if err := RedisClient.Ping(Ctx).Err(); err != nil {
			if os.Getenv("GIN_MODE") == "release" {
				log.Fatal(err)
			}
			fmt.Println(err)
		}
	}
	fmt.Println("redis connection established")
}
func PublishTopic(topic string, payload interface{}) error {
	if err := RedisClient.Publish(Ctx, topic, payload).Err(); err != nil {
		return err
	}
	return nil
}
