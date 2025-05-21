package config

import (
	"github.com/spf13/viper"
)

type StreamConfig struct {
	Name        string
	Subjects    []string
	SubjectName string
	Retention   string
	Storage     string
	MaxAge      int64 // in seconds
}

type Config struct {
	NatsURL        string
	NatsMonitorURL string
	Stream         StreamConfig
}

func Load() (*Config, error) {
	// Set default values
	viper.SetDefault("nats.url", "nats://localhost:4222")
	viper.SetDefault("nats.monitor_url", "http://localhost:8222")
	viper.SetDefault("stream.name", "ORDERS")
	viper.SetDefault("stream.subjects", []string{"ORDERS.*"})
	viper.SetDefault("stream.subjectName", "ORDERS.received")
	viper.SetDefault("stream.retention", "workqueue")
	viper.SetDefault("stream.storage", "file")
	viper.SetDefault("stream.maxAge", 86400) // 24 hours in seconds

	// Set environment variables prefix
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	// Read config file if exists
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	_ = viper.ReadInConfig() // Ignore error if config file does not exist

	// Create config struct
	cfg := &Config{
		NatsURL:        viper.GetString("nats.url"),
		NatsMonitorURL: viper.GetString("nats.monitor_url"),
		Stream: StreamConfig{
			Name:        viper.GetString("stream.name"),
			Subjects:    viper.GetStringSlice("stream.subjects"),
			SubjectName: viper.GetString("stream.subjectName"),
			Retention:   viper.GetString("stream.retention"),
			Storage:     viper.GetString("stream.storage"),
			MaxAge:      viper.GetInt64("stream.maxAge"),
		},
	}

	return cfg, nil
}
