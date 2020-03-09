package handlers

import (
	"github.com/carldanley/homelab-hue/src/cache"
	"github.com/carldanley/homelab-hue/src/events"
	"github.com/sirupsen/logrus"
)

func Register(events events.EventSystem, cache cache.HueCacheSystem, log *logrus.Logger) {
	RegisterHueLightGetHandler(events, cache, log)
	RegisterHueLightSetHandler(events, cache, log)
	RegisterHueSensorGetHandler(events, cache, log)
}
