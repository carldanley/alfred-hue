package events

import (
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type EventSystem struct {
	bus         *nats.Conn
	log         *logrus.Logger
	events      chan Event
	eventPrefix string
}

type Event struct {
	Name        string
	JSONPayload string
}

type RequestHandler interface {
	Process([]byte) (string, error)
}

type RequestHandlerError struct {
	Error string `json:"error"`
}
