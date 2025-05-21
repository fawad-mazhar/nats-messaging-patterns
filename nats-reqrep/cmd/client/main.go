package main

import (
	"context"
	"log"
	"time"

	"github.com/fawadmazhar/nats-reqrep/internal/reqrep"
)

func main() {
	client := reqrep.NewClient()
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.SendRequest(ctx, "Hello from client!")
	if err != nil {
		log.Fatalf("Client error: %v", err)
	}
}
