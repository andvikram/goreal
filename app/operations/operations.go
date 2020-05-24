package op

import (
	"errors"

	"github.com/andvikram/goreal/app/model"
	"github.com/andvikram/goreal/app/operations/pulsar"
	"github.com/andvikram/goreal/ds"
)

// GoRealOp ...
type GoRealOp struct {
	PulsarOp *pulsar.Op
	DSName   string
}

// Operations defines the methods for datasink operations
type Operations interface {
	Send(*model.Message, string) error
	Receive(string, string) (*model.Message, error)
}

var err error

// NewGoRealOp ...
func NewGoRealOp() *GoRealOp {
	gro := new(GoRealOp)
	gro.DSName = ds.DSName
	switch gro.DSName {
	case ds.PulsarDS:
		gro.PulsarOp = pulsar.NewOp()
	}
	return gro
}

// Send ...
func (gorealOp *GoRealOp) Send(message *model.Message) error {
	switch gorealOp.DSName {
	case ds.PulsarDS:
		err = gorealOp.PulsarOp.Send(message)
	default:
		err = errors.New("unsupported sink")
	}
	return err
}

// Receive ...
func (gorealOp *GoRealOp) Receive() (*model.Message, error) {
	var message *model.Message
	switch gorealOp.DSName {
	case ds.PulsarDS:
		message, err = gorealOp.PulsarOp.Receive()
	default:
		err = errors.New("unsupported sink")
	}
	return message, err
}

// InitProducer ...
func (gorealOp *GoRealOp) InitProducer(topicID string) error {
	switch gorealOp.DSName {
	case ds.PulsarDS:
		err = gorealOp.PulsarOp.InitProducer(topicID)
	default:
		err = errors.New("unsupported sink")
	}
	return err
}

// CloseProducer ...
func (gorealOp *GoRealOp) CloseProducer() {
	switch gorealOp.DSName {
	case ds.PulsarDS:
		gorealOp.PulsarOp.CloseProducer()
	default:
		panic("unsupported sink")
	}
}

// InitConsumer ...
func (gorealOp *GoRealOp) InitConsumer(topicID, peer string) error {
	switch gorealOp.DSName {
	case ds.PulsarDS:
		err = gorealOp.PulsarOp.InitConsumer(topicID, peer)
	default:
		err = errors.New("unsupported sink")
	}
	return err
}

// CloseConsumer ...
func (gorealOp *GoRealOp) CloseConsumer() {
	switch gorealOp.DSName {
	case ds.PulsarDS:
		gorealOp.PulsarOp.CloseConsumer()
	default:
		panic("unsupported sink")
	}
}
