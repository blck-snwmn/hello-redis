package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

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

	setAndGet(context.Background(), rdb)
	pubsub(context.Background(), rdb)
}

func setAndGet(ctx context.Context, rdb redis.UniversalClient) {
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
	switch {
	case errors.Is(err, redis.Nil):
		fmt.Println("key2 does not exist")
	case err != nil:
		panic(err)
	default:
		fmt.Println("key2", val2)
	}
}

func pubsub(ctx context.Context, rdb redis.UniversalClient) {
	const cname = "channelname"
	const recieverNum = 1000

	sendMessages := []string{"this is message", "good morning", "foo bar"}

	var sg sync.WaitGroup
	for i := 0; i < recieverNum; i++ {
		sg.Add(1)
		subscriber := rdb.Subscribe(ctx, cname)
		go func(num int, subscriber *redis.PubSub) {
			defer sg.Done()
			for i := 0; i < len(sendMessages); i++ { // recieve 3 message
				msg, err := subscriber.ReceiveMessage(ctx)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("[%3d]channel=%s, payload=%s\n", num, msg.Channel, msg.Payload)
			}
		}(i, subscriber)
	}

	time.Sleep(time.Second) // wait for setup

	// publish 3 message
	for _, msg := range sendMessages {
		if err := rdb.Publish(ctx, cname, msg).Err(); err != nil {
			panic(err)
		}
	}
	sg.Wait()
}
