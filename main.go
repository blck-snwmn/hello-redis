package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func containerAddrToHostAddr(caddr string) string {
	switch caddr {
	case "172.26.0.2:6379":
		return "127.0.0.1:6384"
	case "172.26.0.3:6379":
		return "127.0.0.1:6381"
	case "172.26.0.4:6379":
		return "127.0.0.1:6379"
	case "172.26.0.5:6379":
		return "127.0.0.1:6380"
	case "172.26.0.6:6379":
		return "127.0.0.1:6383"
	case "172.26.0.7:6379":
		return "127.0.0.1:6382"
	default:
		return caddr
	}
}

func main() {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"127.0.0.1:6379",
		},
		NewClient: func(opt *redis.Options) *redis.Client {
			opt.Addr = containerAddrToHostAddr(opt.Addr)
			return redis.NewClient(opt)
		},
	})

	ctx := context.Background()
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if errors.Is(err, redis.Nil) {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}
