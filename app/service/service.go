package service

import (
	"encoding/json"

	"github.com/andvikram/goreal/app/api"
	"github.com/andvikram/goreal/app/model"
	op "github.com/andvikram/goreal/app/operations"
	"github.com/andvikram/goreal/logger"
	"github.com/gorilla/websocket"
)

// GoRealService provides the interfaces for the GoReal service
type GoRealService interface {
	SubscribeTopic(*websocket.Conn)
	PublishToTopic(*api.PublishToTopicRequest) *api.PublishToTopicResponse
}

var (
	// GR is instance for the GoReal service interface
	GR GoRealService
	// Discontinue receiving messages
	Discontinue = false
	log         = logger.GoRealLog{}
	err         error
)

type goReal struct {
	op op.Operations
}

// Initialize ...
func Initialize() {
	GR = newService()
}

func newService() GoRealService {
	gr := new(goReal)
	gr.op = op.NewGoRealOp()
	return gr
}

// SubscribeTopic ...
func (gr *goReal) SubscribeTopic(ws *websocket.Conn) {
	defer ws.Close()

	// Read the initial message to get info
	_, JSONMessage, err := ws.ReadMessage()
	if err != nil {
		log.WithFields(logger.Fields{
			"event": "service.SubscribeTopic()",
			"error": err,
		}).Error("Error reading message from ws client")
		return
	}
	data := make(map[string]string)
	err = json.Unmarshal([]byte(JSONMessage), &data)
	if err != nil {
		log.WithFields(logger.Fields{
			"event": "service.SubscribeTopic()",
			"error": err,
		}).Error("Error unmarshalling JSON message")
		return
	}
	if data["topicID"] == "" || data["peer"] == "" {
		log.WithFields(logger.Fields{
			"event":   "service.SubscribeTopic()",
			"error":   "Parameter values not present",
			"message": data,
		}).Error("Parameter values not present")
		return
	}

	err = gr.op.InitConsumer(data["topicID"], data["peer"])
	defer gr.op.CloseConsumer()
	if err != nil {
		log.WithFields(logger.Fields{
			"event": "service.SubscribeTopic()",
			"error": err,
		}).Error("Error initializing consumer")
		return
	}

	for {
		message, err := gr.op.Receive()
		if err != nil {
			if !Discontinue {
				log.WithFields(logger.Fields{
					"event": "service.SubscribeTopic()",
					"error": err,
				}).Error("Error receiving message")
			}
			return
		}

		err = ws.WriteJSON(message.Data)
		if err != nil {
			log.WithFields(logger.Fields{
				"event": "service.SubscribeTopic()",
				"error": err,
			}).Error("Error writing to WebSocket connection")
			return
		}
	}
}

// PublishToTopic ...
func (gr *goReal) PublishToTopic(req *api.PublishToTopicRequest) *api.PublishToTopicResponse {
	response := new(api.PublishToTopicResponse)

	message := model.NewMessage()
	message.Data = req.Message

	err = gr.op.InitProducer(req.TopicID)
	defer gr.op.CloseProducer()
	if err != nil {
		log.WithFields(logger.Fields{
			"event": "service.PublishToTopic()",
			"error": err,
		}).Error("Error initializing producer")
		response.Status.Code = "500"
		response.Status.Text = "Internal Server Error"
		return response
	}

	err := gr.op.Send(message)
	if err != nil {
		log.WithFields(logger.Fields{
			"event": "service.PublishToTopic()",
			"error": err,
		}).Info("Failed to send message")
		response.Status.Code = "500"
		response.Status.Text = "Internal Server Error"
		return response
	}

	response.Status.Code = "200"
	response.Status.Text = "Success"
	return response
}
