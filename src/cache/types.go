package cache

import (
	"encoding/json"

	"github.com/amimof/huego"
	"github.com/carldanley/alfred-hue/src/events"
	"github.com/sirupsen/logrus"
)

type HueCacheSystem struct {
	bridge *huego.Bridge
	log    *logrus.Logger
	events events.EventSystem

	lights       map[int]HueLight
	sensors      map[int]HueSensor
	shuttingDown bool
}

type HueLight struct {
	ID               int    `json:"id"`
	UniqueID         string `json:"uniqueId"`
	Name             string `json:"name"`
	ModelID          string `json:"modelId"`
	ManufacturerName string `json:"manufacturerName"`
	SwVersion        string `json:"swVersion"`

	On        bool `json:"on"`
	Reachable bool `json:"reachable"`

	Brightness       uint8     `json:"brightness"`
	Hue              uint16    `json:"hue"`
	Saturation       uint8     `json:"saturation"`
	Xy               []float32 `json:"xy"`
	ColorTemperature uint16    `json:"colorTemperature"`
	Effect           string    `json:"effect"`
	TransitionTime   uint16    `json:"transitionTime"`
	ColorMode        string    `json:"colorMode"`
}

func (hl *HueLight) ToJSON() string {
	json, _ := json.Marshal(hl)

	return string(json)
}

type HueSensor struct {
	ID               int    `json:"id"`
	UniqueID         string `json:"uniqueId"`
	Name             string `json:"name"`
	ModelID          string `json:"modelId"`
	ManufacturerName string `json:"manufacturerName"`
	SwVersion        string `json:"swVersion"`

	On        bool    `json:"on"`
	Battery   float64 `json:"battery"`
	Reachable bool    `json:"reachable"`

	LightLevel  float64 `json:"lightLevel"`
	Dark        bool    `json:"dark"`
	Daylight    bool    `json:"daylight"`
	Presence    bool    `json:"presence"`
	Temperature float64 `json:"temperature"`
}

func (hl *HueSensor) ToJSON() string {
	json, _ := json.Marshal(hl)

	return string(json)
}
