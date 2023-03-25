package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
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
	const channelName = "redis-chat"
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"127.0.0.1:6379",
		},
		NewClient: func(opt *redis.Options) *redis.Client {
			opt.Addr = containerAddrToHostAddr(opt.Addr)
			return redis.NewClient(opt)
		},
	})

	scanner := bufio.NewScanner(os.Stdin)

	subscriber := rdb.Subscribe(context.Background(), channelName)
	go func() {
		for msg := range subscriber.Channel() {
			fmt.Printf("[recieve:%v]\n\tmessage: %s\n\tchannel: %s\n", time.Now().Format(time.RFC3339), msg.Payload, msg.Channel)
		}
	}()
	for {
		scanner.Scan()
		msg := scanner.Text()
		switch msg {
		case "", "end":
			return
		}
		err := rdb.Publish(context.Background(), channelName, msg).Err()
		if err != nil {
			log.Fatal(err)
		}
	}
}
