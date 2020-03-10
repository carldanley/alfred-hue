package cache

import (
	"time"

	"github.com/amimof/huego"
	"github.com/carldanley/alfred-hue/src/config"
	"github.com/carldanley/alfred-hue/src/events"
	"github.com/carldanley/alfred-hue/src/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func NewHueCacheSystem(cfg config.Config, log *logrus.Logger, events events.EventSystem) HueCacheSystem {
	system := HueCacheSystem{
		bridge: huego.New(cfg.HueBridgeAddress, cfg.HueUserID),
		log:    log,
		events: events,

		lights:       map[int]HueLight{},
		sensors:      map[int]HueSensor{},
		shuttingDown: false,
	}

	return system
}

func (hcs *HueCacheSystem) Shutdown() {
	hcs.log.Debug("shutting down hue cache system")
	hcs.shuttingDown = true
}

func (hcs *HueCacheSystem) Startup() {
	hcs.log.Debug("starting up hue cache system")

	for {
		if hcs.shuttingDown {
			return
		}

		startTime := time.Now()

		if err := hcs.updateLights(); err != nil {
			metrics.HueCacheUpdateErrorsCounter.With(prometheus.Labels{
				"type": "lights",
			}).Inc()

			hcs.log.WithError(err).Warn("could not update lights")
		}

		if err := hcs.updateSensors(); err != nil {
			metrics.HueCacheUpdateErrorsCounter.With(prometheus.Labels{
				"type": "sensors",
			}).Inc()

			hcs.log.WithError(err).Warn("could not update sensors")
		}

		elapsedTime := time.Since(startTime).Milliseconds()
		metrics.HueCacheUpdateLatencyMSHistogram.Observe(float64(elapsedTime))
		if elapsedTime < CACHE_SYNC_INTERVAL_MS {
			time.Sleep(time.Millisecond * time.Duration(CACHE_SYNC_INTERVAL_MS-elapsedTime))
		}
	}
}
