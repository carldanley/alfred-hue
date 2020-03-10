package handlers

import (
	"encoding/json"

	"github.com/amimof/huego"
	"github.com/carldanley/alfred-hue/src/cache"
	"github.com/carldanley/alfred-hue/src/events"
	"github.com/sirupsen/logrus"
)

type HueLightSetHandler struct {
	cache cache.HueCacheSystem
	log   *logrus.Logger
}

type HueLightSetRequest struct {
	ID    int         `json:"id"`
	State huego.State `json:"state"`
}

type HueLightSetSuccessResponse struct {
	Ok bool `json:"ok"`
}

func (h *HueLightSetHandler) Process(message []byte) (string, error) {
	var request HueLightSetRequest

	if err := json.Unmarshal(message, &request); err != nil {
		return "", err
	}

	err := h.cache.SetLightStateById(request.ID, request.State)
	if err != nil {
		return "", err
	}

	json, _ := json.Marshal(HueLightSetSuccessResponse{Ok: true})
	return string(json), nil
}

func RegisterHueLightSetHandler(events events.EventSystem, cache cache.HueCacheSystem, log *logrus.Logger) {
	handler := HueLightSetHandler{
		cache: cache,
		log:   log,
	}

	events.Subscribe("hue.light.set", &handler)
}
