package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fawadmazhar/nats-pubsub/internal/config"
	"github.com/fawadmazhar/nats-pubsub/internal/pubsub"
	"github.com/fawadmazhar/nats-pubsub/internal/stream"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to NATS
	js, nc, err := stream.Connect(cfg.NatsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create subscriber
	subscriber := pubsub.NewSubscriber(js, cfg.Stream.Name, cfg.Stream.SubjectName)

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start subscribing in a goroutine
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		if err := subscriber.Run(ctx, 10); err != nil { // Max 10 concurrent messages
			log.Printf("Subscriber error: %v", err)
			cancel()
		}
	}()

	// Wait for termination signal
	<-sigCh
	fmt.Println("\nShutting down subscriber...")
	cancel()
	
	// Wait for subscriber to finish processing
	<-doneCh
	fmt.Println("Subscriber shutdown complete")
}
