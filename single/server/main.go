package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:            "localhost:6379",
		PoolSize:        100,
		ConnMaxIdleTime: time.Duration(10) * time.Second,
		ConnMaxLifetime: time.Duration(30) * time.Second,
	})
	// err := client.Ping(context.Background()).Err()
	// if err != nil {
	// 	panic(err)
	// }
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := client.Conn().Set(context.Background(), "key", "value", time.Minute).Err()
		if err != nil {
			log.Println(err)
		} else {
			log.Println("success")
		}
		w.WriteHeader(http.StatusOK)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
