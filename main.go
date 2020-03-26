package main

import (
	"github.com/go-redis/redis/v7"
	"github.com/thanhchungbtc/redis-in-action-go/chapter3"
	"time"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	a := chapter3.NewApp(client)
	//go a.RunPubSub()
	//a.Publisher(6)

	for i := 0; i < 50; i++ {
		go a.NoTrans()
	}

	time.Sleep(time.Second)

}
