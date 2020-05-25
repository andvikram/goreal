package ds

import (
	"github.com/andvikram/goreal/logger"
	"github.com/apache/pulsar-client-go/pulsar"
)

// GoRealDS ...
type GoRealDS struct {
	PulsarDS pulsar.Client
	Env      string
}

const (
	// PulsarDS ...
	PulsarDS = "pulsar"
)

var (
	// DSName ...
	DSName string
	// DSUrl ...
	DSUrl string
	// Env ...
	Env string
	log = logger.GoRealLog{}
	err error
)

// NewGoRealDS ...
func NewGoRealDS() *GoRealDS {
	grds := new(GoRealDS)
	grds.Env = Env
	switch DSName {
	case PulsarDS:
		grds.PulsarDS, err = NewPulsarDS(DSUrl)
	}
	if err != nil {
		log.WithFields(logger.Fields{
			"event": "ds.NewGoRealDS()",
			"error": err,
		}).Fatal("Failed to get DS client")
	}
	return grds
}

// CloseDS ...
func CloseDS() error {
	err = nil
	switch DSName {
	case PulsarDS:
		ClosePulsarDS()
	}
	return err
}
