package main

import (
	"context"
	"github.com/anaregdesign/papaya/concurrent/pubsub"
	"log"
	"runtime"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	topic := pubsub.NewTopic[int]("int_topic")
	sub := topic.NewSubscription("printer", int64(runtime.NumCPU()), time.Minute, time.Minute)

	// Publish 100 messages
	go func(t *pubsub.Topic[int]) {
		for i := 0; i < 100; i++ {
			t.Publish(i)
		}
	}(topic)

	// Subscribe to the topic
	go sub.Subscribe(ctx, func(x *pubsub.Message[int]) {
		// Simulate some work
		time.Sleep(1 * time.Second)
		log.Printf("Value: %v", x.Body())
		x.Ack()
	})

	<-ctx.Done()

}
