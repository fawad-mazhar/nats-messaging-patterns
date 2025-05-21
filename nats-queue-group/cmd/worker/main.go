package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fawadmazhar/nats-queue-group/internal/worker"
)

func main() {
	// Get the worker ID from command line, default to 1
	workerID := "1"
	if len(os.Args) > 1 {
		workerID = os.Args[1]
	}

	// Create and start the worker
	w := worker.NewWorker(workerID)
	
	// Start worker in a goroutine
	go func() {
		if err := w.Run(); err != nil {
			log.Fatalf("Worker error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	
	log.Printf("Worker %s shutting down...\n", workerID)
	w.Shutdown()
}
