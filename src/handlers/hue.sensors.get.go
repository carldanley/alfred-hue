package handlers

import (
	"encoding/json"

	"github.com/carldanley/alfred-hue/src/cache"
	"github.com/carldanley/alfred-hue/src/events"
	"github.com/sirupsen/logrus"
)

type HueSensorsGetHandler struct {
	cache cache.HueCacheSystem
	log   *logrus.Logger
}

func (h *HueSensorsGetHandler) Process(message []byte) (string, error) {
	sensors := h.cache.GetSensors()
	json, _ := json.Marshal(sensors)

	return string(json), nil
}

func RegisterHueSensorsGetHandler(events events.EventSystem, cache cache.HueCacheSystem, log *logrus.Logger) {
	handler := HueSensorsGetHandler{
		cache: cache,
		log:   log,
	}

	events.Subscribe("hue.sensors.get", &handler)
}
