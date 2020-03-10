package main

import (
	"os"
	"time"

	"github.com/carldanley/alfred-hue/src/config"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	requestSubject := "alfred.hue.sensor.get"
	payloadToSend := `{"id": 16,"state": {"on": true,"bri":255,"sat":255,"colormode":"hs","hue":60566,"sat":22,"effect":"none"}}`

	cfg, _ := config.GetConfig()
	conn, err := nats.Connect(
		cfg.NatsServer,
		nats.ReconnectWait(time.Second*10),
		nats.DisconnectErrHandler(onDisconnect),
		nats.ReconnectHandler(onReconnect),
	)

	if err != nil {
		log.WithError(err).Panic("Could not create a NATS connection")
	}

	defer conn.Close()

	msg, err := conn.Request(requestSubject, []byte(payloadToSend), (time.Second * 2))
	if err != nil {
		if conn.LastError() != nil {
			log.WithError(conn.LastError()).Panic("could not send request")
		}
	}

	if msg == nil {
		log.Panic("did not receive response")
	}

	log.Printf("Published [%s] : %s", requestSubject, payloadToSend)
	log.Printf("Received: [%v] : %s", msg.Subject, string(msg.Data))
}

func onDisconnect(conn *nats.Conn, err error) {
	log.Warn("disconnected from event bus")
}

func onReconnect(conn *nats.Conn) {
	log.Info("reconnected to event bus")
}
