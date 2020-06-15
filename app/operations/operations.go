package op

import (
	"errors"

	"github.com/andvikram/goreal/app/model"
	"github.com/andvikram/goreal/app/operations/pulsar"
	"github.com/andvikram/goreal/mb"
)

// Operations defines the methods for message bus operations
type Operations interface {
	InitProducer(string) error
	InitConsumer(string, string) error
	Send(*model.Message) error
	Receive() (*model.Message, error)
	CloseProducer()
	CloseConsumer()
}

type goRealOp struct {
	pulsar *pulsar.Op
	mbName string
}

// NewGoRealOp returns the Operations interface type
func NewGoRealOp() Operations {
	gro := new(goRealOp)
	gro.mbName = mb.MBName
	switch gro.mbName {
	case mb.PulsarMB:
		gro.pulsar = pulsar.NewOp()
	}
	return gro
}

var err error

// Send ...
func (op *goRealOp) Send(message *model.Message) error {
	switch op.mbName {
	case mb.PulsarMB:
		err = op.pulsar.Send(message)
	default:
		err = errors.New("unsupported sink")
	}
	return err
}

// Receive ...
func (op *goRealOp) Receive() (*model.Message, error) {
	var message *model.Message
	switch op.mbName {
	case mb.PulsarMB:
		message, err = op.pulsar.Receive()
	default:
		err = errors.New("unsupported sink")
	}
	return message, err
}

// InitProducer ...
func (op *goRealOp) InitProducer(topicID string) error {
	switch op.mbName {
	case mb.PulsarMB:
		err = op.pulsar.InitProducer(topicID)
	default:
		err = errors.New("unsupported sink")
	}
	return err
}

// CloseProducer ...
func (op *goRealOp) CloseProducer() {
	switch op.mbName {
	case mb.PulsarMB:
		op.pulsar.CloseProducer()
	default:
		panic("unsupported sink")
	}
}

// InitConsumer ...
func (op *goRealOp) InitConsumer(topicID, peerID string) error {
	switch op.mbName {
	case mb.PulsarMB:
		err = op.pulsar.InitConsumer(topicID, peerID)
	default:
		err = errors.New("unsupported sink")
	}
	return err
}

// CloseConsumer ...
func (op *goRealOp) CloseConsumer() {
	switch op.mbName {
	case mb.PulsarMB:
		op.pulsar.CloseConsumer()
	default:
		panic("unsupported sink")
	}
}
