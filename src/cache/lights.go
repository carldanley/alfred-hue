package cache

import (
	"github.com/amimof/huego"
)

func (hcs *HueCacheSystem) updateLights() error {
	lights, err := hcs.bridge.GetLights()
	if err != nil {
		return err
	}

	for _, light := range lights {
		new := hcs.convertHuegoLightToHueLight(light)
		old, err := hcs.GetLightById(light.ID)
		json := new.ToJSON()

		if err != nil {
			hcs.lights[new.ID] = new
			continue
		}

		if new.Name != old.Name {
			hcs.events.Publish("hue.light.name", json)
		}

		if new.SwVersion != old.SwVersion {
			hcs.events.Publish("hue.light.softwareVersion", json)
		}

		if new.On != old.On {
			if new.On {
				hcs.events.Publish("hue.light.on", json)
			} else {
				hcs.events.Publish("hue.light.off", json)
			}
		}

		if new.Reachable != old.Reachable {
			if new.Reachable {
				hcs.events.Publish("hue.light.reachable", json)
			} else {
				hcs.events.Publish("hue.light.unreachable", json)
			}
		}

		if new.Brightness != old.Brightness {
			hcs.events.Publish("hue.light.on", json)
		}

		if new.Hue != old.Hue {
			hcs.events.Publish("hue.light.hue", json)
		}

		if new.Saturation != old.Saturation {
			hcs.events.Publish("hue.light.saturation", json)
		}

		if len(new.Xy) != len(old.Xy) {
			for id, value := range new.Xy {
				if id >= len(old.Xy) {
					hcs.events.Publish("hue.light.xy", json)
					break
				}

				if value != old.Xy[id] {
					hcs.events.Publish("hue.light.xy", json)
					break
				}
			}
		}

		if new.ColorTemperature != old.ColorTemperature {
			hcs.events.Publish("hue.light.colorTemperature", json)
		}

		if new.Effect != old.Effect {
			hcs.events.Publish("hue.light.effect", json)
		}

		if new.TransitionTime != old.TransitionTime {
			hcs.events.Publish("hue.light.transitionTime", json)
		}

		if new.ColorMode != old.ColorMode {
			hcs.events.Publish("hue.light.colorMode", json)
		}

		hcs.lights[new.ID] = new
	}

	return nil
}

func (hcs *HueCacheSystem) convertHuegoLightToHueLight(light huego.Light) HueLight {
	return HueLight{
		ID:               light.ID,
		UniqueID:         light.UniqueID,
		Name:             light.Name,
		ModelID:          light.ModelID,
		ManufacturerName: light.ManufacturerName,
		SwVersion:        light.SwVersion,
		On:               light.State.On,
		Reachable:        light.State.Reachable,
		Brightness:       light.State.Bri,
		Hue:              light.State.Hue,
		Saturation:       light.State.Sat,
		Xy:               light.State.Xy,
		ColorTemperature: light.State.Ct,
		Effect:           light.State.Effect,
		TransitionTime:   light.State.TransitionTime,
		ColorMode:        light.State.ColorMode,
	}
}
