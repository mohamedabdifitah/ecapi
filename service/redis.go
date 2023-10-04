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
func SearchDrivers(limit int, lang, lat, r float64, unit string, withdist bool) []redis.GeoLocation {
	value, err := RedisClient.GeoSearchLocation(Ctx, "driver", &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  lang,
			Latitude:   lat,
			Radius:     r,
			RadiusUnit: unit,
		},
		WithDist: withdist,
	}).Result()
	if err == redis.Nil {
		fmt.Println("key does not exist")
	} else if err != nil {
		panic(err)
	}
	return value
}

// set driver location using redis
func SetDriverLocation(name string, Longitude float64, Latitude float64) (int64, error) {
	res, err := RedisClient.GeoAdd(Ctx, "driver", &redis.GeoLocation{
		Name:      name,
		Longitude: Longitude,
		Latitude:  Latitude,
	}).Result()
	if err != nil {
		return 0, err
	}
	return res, nil
}
