package stream

import (
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/fawadmazhar/nats-pubsub/internal/config"
)

// Connect establishes a connection to NATS and returns JetStream context
func Connect(url string) (nats.JetStreamContext, *nats.Conn, error) {
	// Connect to NATS
	nc, err := nats.Connect(url,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(5),
		nats.ReconnectWait(time.Second),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to NATS: %w", err)
	}

	// Create JetStream Context
	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, nil, fmt.Errorf("error creating JetStream context: %w", err)
	}

	return js, nc, nil
}

// getRetentionPolicy converts string to nats.RetentionPolicy
func getRetentionPolicy(retention string) nats.RetentionPolicy {
	switch strings.ToLower(retention) {
	case "limits":
		return nats.LimitsPolicy
	case "interest":
		return nats.InterestPolicy
	case "workqueue":
		return nats.WorkQueuePolicy
	default:
		return nats.LimitsPolicy
	}
}

// getStorageType converts string to nats.StorageType
func getStorageType(storage string) nats.StorageType {
	switch strings.ToLower(storage) {
	case "file":
		return nats.FileStorage
	case "memory":
		return nats.MemoryStorage
	default:
		return nats.FileStorage
	}
}

// Setup creates the stream if it doesn't exist
func Setup(js nats.JetStreamContext, cfg config.StreamConfig) error {
	// Check if the stream already exists
	stream, err := js.StreamInfo(cfg.Name)
	if err != nil && err != nats.ErrStreamNotFound {
		return fmt.Errorf("error checking stream info: %w", err)
	}

	// If stream doesn't exist, create it
	if stream == nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:       cfg.Name,
			Subjects:   cfg.Subjects,
			Retention:  getRetentionPolicy(cfg.Retention),
			Storage:    getStorageType(cfg.Storage),
			MaxAge:     time.Duration(cfg.MaxAge) * time.Second,
			Replicas:   1,
			Discard:    nats.DiscardOld,
			MaxMsgs:    -1,
			MaxBytes:   -1,
			Duplicates: time.Minute,
		})
		if err != nil {
			return fmt.Errorf("error creating stream: %w", err)
		}
	}

	return nil
}
