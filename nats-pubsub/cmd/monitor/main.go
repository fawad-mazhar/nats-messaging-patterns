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
	"github.com/fawadmazhar/nats-pubsub/internal/monitor"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create monitor service
	monitorService := monitor.NewMonitor(cfg.NatsMonitorURL)

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start monitoring in a goroutine
	go func() {
		// Default to checking every 5 seconds
		if err := monitorService.Run(ctx, 5*time.Second); err != nil {
			log.Printf("Monitor error: %v", err)
			cancel()
		}
	}()

	fmt.Println("NATS monitor started. Press Ctrl+C to exit.")

	// Wait for termination signal
	<-sigCh
	fmt.Println("\nShutting down monitor...")
	cancel()
	
	// Give a moment for any in-flight operations to complete
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Monitor shutdown complete")
}
