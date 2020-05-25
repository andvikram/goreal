package mb

import (
	"github.com/andvikram/goreal/logger"
	"github.com/apache/pulsar-client-go/pulsar"
)

// GoRealMB ...
type GoRealMB struct {
	PulsarMB pulsar.Client
}

const (
	// PulsarMB ...
	PulsarMB = "pulsar"
)

var (
	// MBName ...
	MBName string
	// MBUrl ...
	MBUrl string
	log   = logger.GoRealLog{}
	err   error
)

// NewGoRealMB ...
func NewGoRealMB() *GoRealMB {
	grmb := new(GoRealMB)
	switch MBName {
	case PulsarMB:
		grmb.PulsarMB, err = NewPulsarMB(MBUrl)
	}
	if err != nil {
		log.WithFields(logger.Fields{
			"event": "mb.NewGoRealMB()",
			"error": err,
		}).Fatal("Failed to get MB client")
	}
	return grmb
}

// CloseMB ...
func CloseMB() error {
	err = nil
	switch MBName {
	case PulsarMB:
		ClosePulsarMB()
	}
	return err
}
