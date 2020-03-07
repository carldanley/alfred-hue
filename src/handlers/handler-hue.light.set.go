package handlers

import (
	"encoding/json"

	"github.com/amimof/huego"
	"github.com/carldanley/homelab-hue/src/cache"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type HueLightSetRequest struct {
	ID    int         `json:"id"`
	State huego.State `json:"state"`
}

type HueLightSetErrorResponse struct {
	Error string `json:"error"`
}

type HueLightSetSuccessResponse struct {
	Ok bool `json:"ok"`
}

func HueLightSetHandler(cacheSystem cache.HueCacheSystem, log *logrus.Logger, msg *nats.Msg) {
	log.WithFields(logrus.Fields{
		"subject": msg.Subject,
		"data":    string(msg.Data),
	}).Println("received message")

	var request HueLightSetRequest

	if err := json.Unmarshal(msg.Data, &request); err != nil {
		json, _ := json.Marshal(HueLightGetErrorResponse{
			Error: "invalid JSON",
		})

		msg.Respond(json)
		return
	}

	err := cacheSystem.SetLightStateById(request.ID, request.State)
	if err != nil {
		json, _ := json.Marshal(HueLightGetErrorResponse{
			Error: "could not set light state",
		})

		msg.Respond(json)
		return
	}

	json, _ := json.Marshal(HueLightSetSuccessResponse{Ok: true})
	msg.Respond(json)
}
