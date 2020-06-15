package pulsar

import (
	"context"
	"encoding/json"

	"github.com/andvikram/goreal/app/model"
	"github.com/andvikram/goreal/mb"
	"github.com/apache/pulsar-client-go/pulsar"
)

// Op ...
type Op struct {
	GoRealMB *mb.GoRealMB
}

var (
	producer pulsar.Producer
	consumer pulsar.Consumer
)

// NewOp ...
func NewOp() *Op {
	op := new(Op)
	op.GoRealMB = mb.NewGoRealMB()
	return op
}

// Send ...
func (rOp *Op) Send(message *model.Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: payload,
	})
	return err
}

// Receive ...
func (rOp *Op) Receive() (*model.Message, error) {
	var message model.Message

	msg, err := consumer.Receive(context.Background())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(msg.Payload(), &message)
	if err != nil {
		return nil, err
	}

	consumer.AckID(msg.ID())

	return &message, err
}

// InitProducer ...
func (rOp *Op) InitProducer(topicID string) error {
	var err error
	client := rOp.GoRealMB.PulsarMB
	if producer != nil {
		return nil
	}
	producer, err = client.CreateProducer(pulsar.ProducerOptions{
		Topic: topicID,
	})
	return err
}

// CloseProducer ...
func (rOp *Op) CloseProducer() {
	producer.Close()
}

// InitConsumer ...
func (rOp *Op) InitConsumer(topicID, peerID string) error {
	var err error
	client := rOp.GoRealMB.PulsarMB
	if consumer != nil {
		if consumer.Subscription() == "subcriptionName-"+peerID {
			return nil
		}
		consumer = nil
	}
	consumer, err = client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topicID,
		SubscriptionName: "subcriptionName-" + peerID,
	})
	return err
}

// CloseConsumer ...
func (rOp *Op) CloseConsumer() {
	consumer.Close()
}
