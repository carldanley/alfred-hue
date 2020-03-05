package events

import (
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type EventSystem struct {
	bus    *nats.Conn
	log    *logrus.Logger
	events chan Event
}

type Event struct {
	Name        string
	JSONPayload string
}
