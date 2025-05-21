package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Connect to NATS and setup JetStream
	js, nc, err := stream.Connect(cfg.NatsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Create the stream
	err = stream.Setup(js, cfg.Stream)
	if err != nil {
		log.Fatalf("Failed to setup stream: %v", err)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create publisher
	publisher := pubsub.NewPublisher(js, cfg.Stream.SubjectName)

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start publishing in a goroutine
	go func() {
		if err := publisher.Run(ctx, 2*time.Second); err != nil {
			log.Printf("Publisher error: %v", err)
			cancel()
		}
	}()

	// Wait for termination signal
	<-sigCh
	fmt.Println("\nShutting down publisher...")
	cancel()
	
	// Give a moment for any in-flight operations to complete
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Publisher shutdown complete")
}
