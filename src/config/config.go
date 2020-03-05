package config

import "os"

type Config struct {
	NatsServer       string
	HueBridgeAddress string
	HueUserID        string
}

func GetConfig() (Config, error) {
	cfg := Config{}

	cfg.NatsServer = os.Getenv("NATS_SERVER")
	if cfg.NatsServer == "" {
		cfg.NatsServer = "nats://127.0.0.1:4222"
	}

	cfg.HueBridgeAddress = os.Getenv("HUE_BRIDGE_ADDRESS")
	cfg.HueUserID = os.Getenv("HUE_USER_ID")

	return cfg, nil
}
