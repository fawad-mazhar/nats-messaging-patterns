package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fawadmazhar/nats-reqrep/internal/reqrep"
)

func main() {
	server := reqrep.NewServer()

	// Start server in a goroutine
	go func() {
		if err := server.Run(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	log.Println("Server is running. Waiting for requests...")

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	
	log.Println("Server shutting down...")
}
