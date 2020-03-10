package config

import (
	"os"
	"strconv"
)

type Config struct {
	NatsServer       string
	NatsEventPrefix  string
	NatsQueueName    string
	HueBridgeAddress string
	HueUserID        string
	MetricsPort      int
}

func GetConfig() (Config, error) {
	cfg := Config{
		NatsServer:       os.Getenv("NATS_SERVER"),
		NatsEventPrefix:  os.Getenv("NATS_EVENT_PREFIX"),
		NatsQueueName:    os.Getenv("NATS_QUEUE_NAME"),
		HueBridgeAddress: os.Getenv("HUE_BRIDGE_ADDRESS"),
		HueUserID:        os.Getenv("HUE_USER_ID"),
		MetricsPort:      9200,
	}

	if cfg.NatsServer == "" {
		cfg.NatsServer = "nats://127.0.0.1:4222"
	}

	if cfg.NatsEventPrefix == "" {
		cfg.NatsEventPrefix = "alfred"
	}

	if cfg.NatsQueueName == "" {
		cfg.NatsQueueName = "alfred-hue"
	}

	metricsPort, err := strconv.Atoi(os.Getenv("METRICS_PORT"))
	if err == nil {
		cfg.MetricsPort = metricsPort
	}

	return cfg, nil
}
