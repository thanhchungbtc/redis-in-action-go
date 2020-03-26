package chapter3

import (
	"log"
	"strconv"
	"sync/atomic"
	"time"
)

func (a *App) Publisher(n int) {
	time.Sleep(time.Second)

	for i := 0; i < n; i++ {
		a.conn.Publish("channel", "message"+strconv.Itoa(i))
		time.Sleep(time.Second)
	}
}

func (a *App) RunPubSub() {
	pubsub := a.conn.Subscribe("channel")
	defer pubsub.Close()

	var count int32 = 0
	for item := range pubsub.Channel() {
		log.Println(item)
		atomic.AddInt32(&count, 1)
		log.Println(count)

		if count == 4 {
			if err := pubsub.Unsubscribe("channel"); err != nil {
				log.Fatalf("%+v")
			} else {
				log.Println("Unsuscribed!")
			}
		}
		if count == 5 {
			break
		}
	}
	log.Println("Done")
}
