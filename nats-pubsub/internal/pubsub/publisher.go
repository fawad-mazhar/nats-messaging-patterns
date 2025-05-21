package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// Publisher handles publishing messages to NATS JetStream
type Publisher struct {
	js          nats.JetStreamContext
	subjectName string
}

// Message represents a sample message structure
type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// NewPublisher creates a new publisher instance
func NewPublisher(js nats.JetStreamContext, subjectName string) *Publisher {
	return &Publisher{
		js:          js,
		subjectName: subjectName,
	}
}

// Run starts the publishing process with the specified interval
func (p *Publisher) Run(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	msgCount := 0

	for {
		select {
		case <-ticker.C:
			msgCount++
			msg := Message{
				ID:        fmt.Sprintf("msg-%d", msgCount),
				Content:   fmt.Sprintf("Message content %d", msgCount),
				Timestamp: time.Now(),
			}

			data, err := json.Marshal(msg)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}

			// Publish message with message ID
			_, err = p.js.Publish(p.subjectName, data, nats.MsgId(msg.ID))
			if err != nil {
				log.Printf("Error publishing message: %v", err)
				continue
			}

			log.Printf("Published message: %s", msg.ID)

		case <-ctx.Done():
			return nil
		}
	}
}
