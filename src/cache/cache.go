package cache

import (
	"errors"
	"time"

	"github.com/amimof/huego"
	"github.com/carldanley/homelab-hue/src/config"
	"github.com/carldanley/homelab-hue/src/events"
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
			hcs.log.WithError(err).Warn("could not update lights")
		}

		if err := hcs.updateSensors(); err != nil {
			hcs.log.WithError(err).Warn("could not update sensors")
		}

		elapsedTime := time.Since(startTime).Milliseconds()
		if elapsedTime < CACHE_SYNC_INTERVAL_MS {
			time.Sleep(time.Millisecond * time.Duration(CACHE_SYNC_INTERVAL_MS-elapsedTime))
		}
	}
}

func (hcs *HueCacheSystem) GetLightById(id int) (HueLight, error) {
	light, ok := hcs.lights[id]

	if !ok {
		return HueLight{}, errors.New("light does not exist in cache")
	}

	return light, nil
}

func (hcs *HueCacheSystem) GetSensorById(id int) (HueSensor, error) {
	sensor, ok := hcs.sensors[id]

	if !ok {
		return HueSensor{}, errors.New("sensor does not exist in cache")
	}

	return sensor, nil
}
