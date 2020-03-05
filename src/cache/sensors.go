package cache

import (
	"github.com/amimof/huego"
)

func (hcs *HueCacheSystem) updateSensors() error {
	sensors, err := hcs.bridge.GetSensors()
	if err != nil {
		return err
	}

	for _, sensor := range sensors {
		if sensor.Type == SENSOR_TYPE_DAYLIGHT {
			continue
		}

		new := hcs.convertHuegoSensorToHueSensor(sensor)
		old, err := hcs.GetSensorById(sensor.ID)
		json := new.ToJSON()

		if err != nil {
			hcs.sensors[new.ID] = new
			continue
		}

		if new.Name != old.Name {
			hcs.events.Publish("hue.sensor.name", json)
		}

		if new.SwVersion != old.SwVersion {
			hcs.events.Publish("hue.sensor.softwareVersion", json)
		}

		if new.On != old.On {
			if new.On {
				hcs.events.Publish("hue.sensor.on", json)
			} else {
				hcs.events.Publish("hue.sensor.off", json)
			}
		}

		if new.Battery != old.Battery {
			hcs.events.Publish("hue.sensor.battery", json)
		}

		if new.Reachable != old.Reachable {
			if new.Reachable {
				hcs.events.Publish("hue.sensor.reachable", json)
			} else {
				hcs.events.Publish("hue.sensor.unreachable", json)
			}
		}

		switch sensor.Type {
		case SENSOR_TYPE_DAYLIGHT:
			break
		case SENSOR_TYPE_LIGHT_LEVEL:

			if new.LightLevel != old.LightLevel {
				hcs.events.Publish("hue.sensor.lightLevel", json)
			}

			if new.Dark != old.Dark {
				hcs.events.Publish("hue.sensor.dark", json)
			}

			if new.Daylight != old.Daylight {
				hcs.events.Publish("hue.sensor.daylight", json)
			}

		case SENSOR_TYPE_PRESENCE:

			if new.Presence != old.Presence {
				hcs.events.Publish("hue.sensor.presence", json)
			}

		case SENSOR_TYPE_TEMPERATURE:

			if new.Temperature != old.Temperature {
				hcs.events.Publish("hue.sensor.temperature", json)
			}

		}

		hcs.sensors[new.ID] = new
	}

	return nil
}

func (hcs *HueCacheSystem) convertHuegoSensorToHueSensor(sensor huego.Sensor) HueSensor {
	converted := HueSensor{
		ID:               sensor.ID,
		UniqueID:         sensor.UniqueID,
		Name:             sensor.Name,
		ModelID:          sensor.ModelID,
		ManufacturerName: sensor.ManufacturerName,
		SwVersion:        sensor.SwVersion,
		On:               sensor.Config["on"].(bool),
	}

	switch sensor.Type {
	case SENSOR_TYPE_DAYLIGHT:
		return converted
	case SENSOR_TYPE_LIGHT_LEVEL:
		converted.LightLevel = hcs.float64(sensor.State["lightlevel"])
		converted.Dark = hcs.bool(sensor.State["dark"])
		converted.Daylight = hcs.bool(sensor.State["daylight"])
	case SENSOR_TYPE_PRESENCE:
		converted.Presence = hcs.bool(sensor.State["presence"])
	case SENSOR_TYPE_TEMPERATURE:
		converted.Temperature = hcs.float64(sensor.State["temperature"])
		converted.Temperature = hcs.convertTemperatureToFahrenheit(converted.Temperature)
	}

	converted.Battery = hcs.float64(sensor.Config["battery"])
	converted.Reachable = hcs.bool(sensor.Config["reachable"])

	return converted
}

func (hcs *HueCacheSystem) float64(value interface{}) float64 {
	if value == nil {
		return float64(0)
	}

	return value.(float64)
}

func (hcs *HueCacheSystem) bool(value interface{}) bool {
	if value == nil {
		return false
	}

	return value.(bool)
}

func (hcs *HueCacheSystem) convertTemperatureToFahrenheit(temperature float64) float64 {
	return ((temperature / 100.0) * 1.8) + 32.0
}
