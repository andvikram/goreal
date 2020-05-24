package pulsar

import (
	"context"
	"encoding/json"

	"github.com/andvikram/goreal/app/model"
	"github.com/andvikram/goreal/ds"
	"github.com/apache/pulsar-client-go/pulsar"
)

// Op ...
type Op struct {
	GoRealDS *ds.GoRealDS
}

var (
	producer pulsar.Producer
	consumer pulsar.Consumer
)

// NewOp ...
func NewOp() *Op {
	op := new(Op)
	op.GoRealDS = ds.NewGoRealDS()
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
	client := rOp.GoRealDS.PulsarDS
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
func (rOp *Op) InitConsumer(topicID, peer string) error {
	var err error
	client := rOp.GoRealDS.PulsarDS
	if consumer != nil {
		if consumer.Subscription() == "subcriptionName-"+peer {
			return nil
		}
		consumer = nil
	}
	consumer, err = client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topicID,
		SubscriptionName: "subcriptionName-" + peer,
		Type:             pulsar.Shared,
	})
	return err
}

// CloseConsumer ...
func (rOp *Op) CloseConsumer() {
	consumer.Close()
}
