package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/rueidis"
)

func main() {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"127.0.0.1:6379"},
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	{
		cmd := client.B().Hset().
			Key("hkey").
			FieldValue().
			FieldValue("innserkey", "xxxxxxxxxxxxxxxx")

		result := client.Do(context.Background(), cmd.Build())
		if result.Error() != nil {
			log.Printf("failed: %s\n", result.Error())
			return
		}
	}
	{
		cmd := client.B().Hget().Key("hkey").Field("innserkey")
		result := client.Do(context.Background(), cmd.Build())
		if result.Error() != nil {
			log.Printf("failed: %s\n", result.Error())
		}
		v, err := result.ToString()
		fmt.Printf("value=`%v`, err=`%v`\n", v, err)
	}
	{
		cmd := client.B().Hget().Key("hkey").Field("otherkey")
		result := client.Do(context.Background(), cmd.Build())
		if result.Error() != nil && !rueidis.IsRedisNil(result.Error()) {
			log.Printf("failed: %s\n", result.Error())
		}
		v, err := result.ToString()
		fmt.Printf("value=`%v`, err=`%v`\n", v, err)
	}
}
