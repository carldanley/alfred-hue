package handlers

import (
	"encoding/json"

	"github.com/carldanley/homelab-hue/src/cache"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type HueSensorGetRequest struct {
	ID int `json:"id"`
}

type HueSensorGetErrorResponse struct {
	Error string `json:"error"`
}

func HueSensorGetHandler(cacheSystem cache.HueCacheSystem, log *logrus.Logger, msg *nats.Msg) {
	log.WithFields(logrus.Fields{
		"subject": msg.Subject,
		"data":    string(msg.Data),
	}).Println("received message")

	var request HueSensorGetRequest

	if err := json.Unmarshal(msg.Data, &request); err != nil {
		json, _ := json.Marshal(HueSensorGetErrorResponse{
			Error: "invalid JSON",
		})

		msg.Respond(json)
		return
	}

	sensor, err := cacheSystem.GetSensorById(request.ID)
	if err != nil {
		json, _ := json.Marshal(HueSensorGetErrorResponse{
			Error: err.Error(),
		})

		msg.Respond(json)
		return
	}

	json, _ := json.Marshal(sensor)
	msg.Respond(json)
}
