package reqrep

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/fawadmazhar/nats-reqrep/internal/config"
)

type Server struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func NewServer() *Server {
	nc, err := nats.Connect(config.DefaultNatsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		log.Fatalf("Failed to get JetStream context: %v", err)
	}

	return &Server{
		nc: nc,
		js: js,
	}
}

func (s *Server) Run() error {
	defer s.nc.Close()

	// Create a stream if it doesn't exist
	_, err := s.js.AddStream(&nats.StreamConfig{
		Name:     config.StreamName,
		Subjects: []string{config.Subject + ".*"},
		Storage:  nats.FileStorage,
	})
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}

	// Create a durable consumer
	_, err = s.js.AddConsumer(config.StreamName, &nats.ConsumerConfig{
		Durable:       config.ConsumerName,
		AckPolicy:     nats.AckExplicitPolicy,
		FilterSubject: config.Subject + ".*",
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	// Subscribe to requests
	sub, err := s.js.PullSubscribe(config.Subject+".*", config.ConsumerName)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	// Process messages
	for {
		msgs, err := sub.Fetch(1, nats.MaxWait(5*time.Second))
		if err != nil {
			if err == nats.ErrTimeout {
				continue
			}
			log.Printf("Error fetching messages: %v", err)
			continue
		}

		for _, msg := range msgs {
			// Process request
			log.Printf("Received request: %s", string(msg.Data))

			// Create response
			response := fmt.Sprintf("Response to: %s", string(msg.Data))

			// Reply to the message
			err = msg.Respond([]byte(response))
			if err != nil {
				log.Printf("Error responding to message: %v", err)
			} else {
				log.Printf("Sent response: %s", response)
			}

			// Acknowledge message
			msg.Ack()
		}
	}
}
