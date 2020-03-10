package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/carldanley/alfred-hue/src/cache"
	"github.com/carldanley/alfred-hue/src/config"
	"github.com/carldanley/alfred-hue/src/events"
	"github.com/carldanley/alfred-hue/src/handlers"
	"github.com/carldanley/alfred-hue/src/metrics"
	"github.com/sirupsen/logrus"
)

var signalChannel chan os.Signal
var log *logrus.Logger

func init() {
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT)

	logLevel := flag.Int("v", 0, "the level of verbosity for logging")
	flag.Parse()

	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)

	switch *logLevel {
	case 0:
		log.SetLevel(logrus.ErrorLevel)
	case 1:
		log.SetLevel(logrus.WarnLevel)
	case 2:
		log.SetLevel(logrus.InfoLevel)
	default:
		log.SetLevel(logrus.DebugLevel)
	}
}

func main() {
	// get the configuration
	cfg, _ := config.GetConfig()

	// create a new event system
	eventSystem := events.NewEventSystem(cfg, log)
	defer eventSystem.Shutdown()

	// create a new hue cache system
	hueCacheSystem := cache.NewHueCacheSystem(cfg, log, eventSystem)
	defer hueCacheSystem.Shutdown()

	// register all of our HUE service subscriptions
	handlers.Register(eventSystem, hueCacheSystem, log)

	// begin processing events
	go metrics.Startup(cfg.MetricsPort, log)
	go eventSystem.Startup()
	go hueCacheSystem.Startup()

	// wait for an interrupt signal
	<-signalChannel

	// drain the event system
	eventSystem.Drain()
}
