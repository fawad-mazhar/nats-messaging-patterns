package config

import "github.com/nats-io/nats.go"

const (
	Subject     = "service.request"
	StreamName  = "REQUESTS"
	ConsumerName = "SERVER"
)

var (
	DefaultNatsURL = nats.DefaultURL
)
