package main

import (
	"context"
	"github.com/anaregdesign/papaya/cache"
	"log"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	c := cache.NewCache[string, string](ctx, 1*time.Second)

	// Set a value
	c.Set("key", "value")

	// Get a value
	if value, ok := c.Get("key"); ok {
		log.Printf("value: %v", value)
	}

	// Wait for the cache to expire
	time.Sleep(1 * time.Second)

	// Get a value
	if _, ok := c.Get("key"); !ok {
		log.Printf("Key expired")
	}

	<-ctx.Done()

}
