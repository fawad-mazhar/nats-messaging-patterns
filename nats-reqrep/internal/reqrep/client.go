package reqrep

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/fawadmazhar/nats-reqrep/internal/config"
)

type Client struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func NewClient() *Client {
	nc, err := nats.Connect(config.DefaultNatsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		log.Fatalf("Failed to get JetStream context: %v", err)
	}

	return &Client{
		nc: nc,
		js: js,
	}
}

func (c *Client) SendRequest(ctx context.Context, data string) error {
	defer c.nc.Close()

	// Create unique subject for this request
	requestSubject := fmt.Sprintf("%s.%d", config.Subject, time.Now().UnixNano())
	
	// Send request using core NATS request-reply
	msg, err := c.nc.RequestWithContext(ctx, requestSubject, []byte(data))
	if err != nil {
		return fmt.Errorf("failed to receive response: %w", err)
	}

	log.Printf("Received response: %s", string(msg.Data))
	return nil
}
