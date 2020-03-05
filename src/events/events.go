package events

import (
	"regexp"
	"time"

	"github.com/carldanley/homelab-hue/src/config"
	"github.com/carldanley/homelab-hue/src/metrics"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func NewEventSystem(cfg config.Config, log *logrus.Logger) EventSystem {
	system := EventSystem{
		log:    log,
		events: make(chan Event),
	}

	conn, err := nats.Connect(
		cfg.NatsServer,
		nats.ReconnectWait(time.Second*10),
		nats.DisconnectErrHandler(system.onDisconnect),
		nats.ReconnectHandler(system.onReconnect),
	)

	if err != nil {
		log.WithError(err).Panic("Could not create a NATS connection")
	}

	system.bus = conn
	return system
}

func (es *EventSystem) onDisconnect(conn *nats.Conn, err error) {
	es.log.Warn("disconnected from event bus")
}

func (es *EventSystem) onReconnect(conn *nats.Conn) {
	es.log.Info("reconnected to event bus")
}

func (es *EventSystem) Shutdown() {
	es.log.Debug("shutting down event system")

	// close the connection to our event bus
	es.bus.Close()

	// close the event stream
	close(es.events)
}

func (es *EventSystem) Startup() {
	es.log.Debug("starting up event system")

	for event := range es.events {
		if !es.bus.IsConnected() {
			es.log.WithField("name", event.Name).Println("skipping event; not connected")
			continue
		}

		es.log.WithFields(logrus.Fields{
			"name":    event.Name,
			"payload": event.JSONPayload,
		}).Debug("processing event")

		if err := es.bus.Publish(event.Name, []byte(event.JSONPayload)); err != nil {
			es.log.WithError(err).Warn("could not publish event to event bus")
		} else {
			isLight, _ := regexp.Match(`^hue.light.*`, []byte(event.Name))
			isSensor, _ := regexp.Match(`^hue.sensor.*`, []byte(event.Name))

			if isLight {
				metrics.HueEventsEmittedCounter.With(prometheus.Labels{
					"event": event.Name,
					"type":  "light",
				}).Inc()
			} else if isSensor {
				metrics.HueEventsEmittedCounter.With(prometheus.Labels{
					"event": event.Name,
					"type":  "sensor",
				}).Inc()
			}
		}
	}
}

func (es *EventSystem) Publish(name, jsonPayload string) {
	es.events <- Event{
		Name:        name,
		JSONPayload: jsonPayload,
	}
}
