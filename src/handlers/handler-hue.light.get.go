package handlers

import (
	"encoding/json"

	"github.com/carldanley/homelab-hue/src/cache"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type HueLightGetRequest struct {
	ID int `json:"id"`
}

type HueLightGetErrorResponse struct {
	Error string `json:"error"`
}

func HueLightGetHandler(cacheSystem cache.HueCacheSystem, log *logrus.Logger, msg *nats.Msg) {
	log.WithFields(logrus.Fields{
		"subject": msg.Subject,
		"data":    string(msg.Data),
	}).Println("received message")

	var request HueLightGetRequest

	if err := json.Unmarshal(msg.Data, &request); err != nil {
		json, _ := json.Marshal(HueLightGetErrorResponse{
			Error: "invalid JSON",
		})

		msg.Respond(json)
		return
	}

	light, err := cacheSystem.GetLightById(request.ID)
	if err != nil {
		json, _ := json.Marshal(HueLightGetErrorResponse{
			Error: err.Error(),
		})

		msg.Respond(json)
		return
	}

	json, _ := json.Marshal(light)
	msg.Respond(json)
}
