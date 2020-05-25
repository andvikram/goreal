package controller

import (
	"encoding/json"
	"net/http"

	"github.com/andvikram/goreal/app/api"
	"github.com/andvikram/goreal/app/service"
	"github.com/andvikram/goreal/configuration"
	"github.com/andvikram/goreal/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// SubscribeTopic opens a websocket connection for the peer
func SubscribeTopic(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to websocket
	upgrader := websocket.Upgrader{
		ReadBufferSize:  50000,
		WriteBufferSize: 50000,
		CheckOrigin: func(r *http.Request) bool {
			if contains(r.Header["Origin"], configuration.Config.Origin) {
				return true
			}
			return false
		},
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithFields(logger.Fields{
			"event": "controller.SubscribeTopic()",
			"error": err,
		}).Error("Failed to upgrade connection to WebSocket")
		return
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	service.GR.SubscribeTopic(ws)
}

// PublishToTopic publishes message to the topic
func PublishToTopic(w http.ResponseWriter, r *http.Request) {
	response := new(api.PublishToTopicResponse)
	publishToTopicReq := new(api.PublishToTopicRequest)
	params := mux.Vars(r)

	var body map[string]string
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.WithFields(logger.Fields{
			"event": "controller.PublishToTopic()",
			"error": err,
		}).Error("Failed to unmarshal request body")
		response.Status.Code = "500"
		response.Status.Text = "Internal Server Error"
		json.NewEncoder(w).Encode(response)
		return
	}
	publishToTopicReq.TopicID = params["topic_id"]
	publishToTopicReq.Message = body["message"]

	response = service.GR.PublishToTopic(publishToTopicReq)
	json.NewEncoder(w).Encode(response)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
