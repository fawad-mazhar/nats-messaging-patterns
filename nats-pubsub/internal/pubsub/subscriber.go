package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
)

// Subscriber handles consuming messages from NATS JetStream
type Subscriber struct {
	js          nats.JetStreamContext
	streamName  string
	subjectName string
}

// NewSubscriber creates a new subscriber instance
func NewSubscriber(js nats.JetStreamContext, streamName, subjectName string) *Subscriber {
	return &Subscriber{
		js:          js,
		streamName:  streamName,
		subjectName: subjectName,
	}
}

// Run starts the subscription process with the specified worker count
func (s *Subscriber) Run(ctx context.Context, maxWorkers int) error {
	// Create a pull subscription
	sub, err := s.js.PullSubscribe(
		s.subjectName,
		fmt.Sprintf("%s-consumer", s.streamName),
		nats.BindStream(s.streamName),
	)
	if err != nil {
		return fmt.Errorf("error creating subscription: %w", err)
	}
	defer sub.Unsubscribe()

	// Create a worker pool
	var wg sync.WaitGroup
	workChan := make(chan *nats.Msg, maxWorkers)

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			s.worker(ctx, workerID, workChan)
		}(i + 1)
	}

	// Main loop for fetching messages
	for {
		select {
		case <-ctx.Done():
			close(workChan)
			wg.Wait()
			return nil
		default:
			msgs, err := sub.Fetch(10, nats.Context(ctx))
			if err != nil {
				if err == context.Canceled {
					continue
				}
				log.Printf("Error fetching messages: %v", err)
				continue
			}

			for _, msg := range msgs {
				select {
				case <-ctx.Done():
					return nil
				case workChan <- msg:
				}
			}
		}
	}
}

// worker processes messages from the work channel
func (s *Subscriber) worker(ctx context.Context, id int, workChan <-chan *nats.Msg) {
	log.Printf("Worker %d started", id)
	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d shutting down", id)
			return
		case msg, ok := <-workChan:
			if !ok {
				log.Printf("Worker %d channel closed", id)
				return
			}

			// Process the message
			var message Message
			if err := json.Unmarshal(msg.Data, &message); err != nil {
				log.Printf("Worker %d - Error unmarshaling message: %v", id, err)
				msg.Nak()
				continue
			}

			// Log message details
			log.Printf("Worker %d - Received message: %s, Content: %s, Timestamp: %v",
				id, message.ID, message.Content, message.Timestamp)

			// Simulate some processing time
			// time.Sleep(100 * time.Millisecond)

			// Acknowledge the message
			if err := msg.Ack(); err != nil {
				log.Printf("Worker %d - Error acknowledging message: %v", id, err)
			}
		}
	}
}
