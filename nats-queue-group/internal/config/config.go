package config

import (
	"time"

	"github.com/nats-io/nats.go"
)

const (
	Subject    = "tasks"
	QueueGroup = "workers"
	NumTasks   = 10
)

var (
	DefaultNatsURL = nats.DefaultURL

	// Task configuration
	TaskProcessingMinTime = 100 * time.Millisecond
	TaskProcessingMaxTime = 500 * time.Millisecond
	TaskPublishInterval  = 100 * time.Millisecond
	TaskResponseTimeout  = 2 * time.Second
)
