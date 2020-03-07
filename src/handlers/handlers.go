package handlers

import (
	"github.com/carldanley/homelab-hue/src/cache"
	"github.com/carldanley/homelab-hue/src/events"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type EventSystemHandler func(cache.HueCacheSystem, *logrus.Logger, *nats.Msg)
type EventSystemHandlerWrapper func(handler EventSystemHandler) nats.MsgHandler

func Register(eventSystem events.EventSystem, cacheSystem cache.HueCacheSystem, log *logrus.Logger) {
	wrap := Wrap(cacheSystem, log)

	eventSystem.Subscribe("hue.light.get", wrap(HueLightGetHandler))
	eventSystem.Subscribe("hue.light.set", wrap(HueLightSetHandler))
	eventSystem.Subscribe("hue.sensor.get", wrap(HueSensorGetHandler))
}

func Wrap(cacheSystem cache.HueCacheSystem, log *logrus.Logger) EventSystemHandlerWrapper {
	return func(handler EventSystemHandler) nats.MsgHandler {
		return func(msg *nats.Msg) {
			handler(cacheSystem, log, msg)
		}
	}
}
