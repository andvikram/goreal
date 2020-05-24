package routes

import (
	"fmt"
	"sync"

	"github.com/andvikram/goreal/app/controller"
	"github.com/andvikram/goreal/configuration"
	"github.com/gorilla/mux"
)

var (
	router *mux.Router
	once   sync.Once
	err    error

	// Scheme ...
	Scheme string
	// Host ...
	Host string
	// Port ...
	Port string
	// BaseURL provides the complete host URL string
	BaseURL string
)

const (
	// SubscribeTopicRoute defines Websocket path to subscribe and listen to a topic
	SubscribeTopicRoute = "/subscribe"
	// PublishToTopicRoute defines API path for publishing a message in a topic
	PublishToTopicRoute = "/topics/{topic_id}/messages/new"
)

// ServiceRoutes defines the API routes and their handlers for the service
func ServiceRoutes() *mux.Router {
	initialize()
	router := mux.NewRouter()

	router.HandleFunc(SubscribeTopicRoute, controller.SubscribeTopic)
	router.HandleFunc(PublishToTopicRoute, controller.PublishToTopic).Methods("POST")

	return router
}

func initialize() {
	config := configuration.Config
	Scheme = config.AppScheme
	Host = config.AppHost
	Port = config.AppPort
	BaseURL = fmt.Sprintf("%s://%s%s", config.AppScheme, config.AppHost, config.AppPort)
}

func init() {
	controller.RouteMap["SubscribeTopicRoute"] = SubscribeTopicRoute
	controller.RouteMap["PublishToTopicRoute"] = PublishToTopicRoute
}
