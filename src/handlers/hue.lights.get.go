package handlers

import (
	"encoding/json"

	"github.com/carldanley/alfred-hue/src/cache"
	"github.com/carldanley/alfred-hue/src/events"
	"github.com/sirupsen/logrus"
)

type HueLightsGetHandler struct {
	cache cache.HueCacheSystem
	log   *logrus.Logger
}

func (h *HueLightsGetHandler) Process(message []byte) (string, error) {
	lights := h.cache.GetLights()
	json, _ := json.Marshal(lights)

	return string(json), nil
}

func RegisterHueLightsGetHandler(events events.EventSystem, cache cache.HueCacheSystem, log *logrus.Logger) {
	handler := HueLightsGetHandler{
		cache: cache,
		log:   log,
	}

	events.Subscribe("hue.lights.get", &handler)
}
