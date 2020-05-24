package api

// SubscribeTopicRequest defines the request struct for SubscribeTopic endpoint
type SubscribeTopicRequest struct {
	PeerID  string `json:"peer_id"`
	TopicID string `json:"topic_id"`
}

// PublishToTopicRequest defines the request struct for PublishToTopic endpoint
type PublishToTopicRequest struct {
	TopicID  string            `json:"topic_id"`
	Message  string            `json:"message"`
	MetaData map[string]string `json:"meta_data"`
}

// PublishToTopicResponse defines the response for PublishToTopic endpoint
type PublishToTopicResponse struct {
	Status
}
