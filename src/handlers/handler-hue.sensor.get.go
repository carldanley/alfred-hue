package handlers

import (
	"encoding/json"

	"github.com/carldanley/alfred-hue/src/cache"
	"github.com/carldanley/alfred-hue/src/events"
	"github.com/sirupsen/logrus"
)

type HueSensorGetHandler struct {
	cache cache.HueCacheSystem
	log   *logrus.Logger
}

type HueSensorGetRequest struct {
	ID int `json:"id"`
}

func (h *HueSensorGetHandler) Process(message []byte) (string, error) {
	var request HueSensorGetRequest

	if err := json.Unmarshal(message, &request); err != nil {
		return "", err
	}

	sensor, err := h.cache.GetSensorById(request.ID)
	if err != nil {
		return "", err
	}

	json, _ := json.Marshal(sensor)
	return string(json), nil
}

func RegisterHueSensorGetHandler(events events.EventSystem, cache cache.HueCacheSystem, log *logrus.Logger) {
	handler := HueSensorGetHandler{
		cache: cache,
		log:   log,
	}

	events.Subscribe("hue.sensor.get", &handler)
}
