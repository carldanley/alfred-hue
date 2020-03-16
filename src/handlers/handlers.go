package handlers

import (
	"github.com/carldanley/alfred-hue/src/cache"
	"github.com/carldanley/alfred-hue/src/events"
	"github.com/sirupsen/logrus"
)

func Register(events events.EventSystem, cache cache.HueCacheSystem, log *logrus.Logger) {
	// lights
	RegisterHueLightGetHandler(events, cache, log)
	RegisterHueLightSetHandler(events, cache, log)
	RegisterHueLightsGetHandler(events, cache, log)

	// sensors
	RegisterHueSensorGetHandler(events, cache, log)
	RegisterHueSensorsGetHandler(events, cache, log)
}
