package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fawadmazhar/nats-queue-group/internal/publisher"
)

func main() {
	// Create and start the publisher
	p := publisher.NewPublisher()
	
	// Start publisher in a goroutine
	go func() {
		if err := p.Run(); err != nil {
			log.Fatalf("Publisher error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	
	log.Println("Publisher shutting down...")
	p.Shutdown()
}

