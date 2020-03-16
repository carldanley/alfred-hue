package handlers

import (
	"encoding/json"

	"github.com/carldanley/alfred-hue/src/cache"
	"github.com/carldanley/alfred-hue/src/events"
	"github.com/sirupsen/logrus"
)

type HueLightGetHandler struct {
	cache cache.HueCacheSystem
	log   *logrus.Logger
}

type HueLightGetRequest struct {
	ID string `json:"id"`
}

func (h *HueLightGetHandler) Process(message []byte) (string, error) {
	var request HueLightGetRequest

	if err := json.Unmarshal(message, &request); err != nil {
		return "", err
	}

	light, err := h.cache.GetLightById(request.ID)
	if err != nil {
		return "", err
	}

	json, _ := json.Marshal(light)
	return string(json), nil
}

func RegisterHueLightGetHandler(events events.EventSystem, cache cache.HueCacheSystem, log *logrus.Logger) {
	handler := HueLightGetHandler{
		cache: cache,
		log:   log,
	}

	events.Subscribe("hue.light.get", &handler)
}
