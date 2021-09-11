package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "SUIbianla123",
		DB:       0,
	})
	// err := rdb.Set(ctx, "ca", "niao", 0).Err()
	// if err != nil {
	// 	fmt.Printf("%v", err)
	// 	return
	// }
	value, err := rdb.Get(ctx, "ca").Result()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%v", value)
}
