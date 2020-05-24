package ds

import (
	"sync"

	"github.com/apache/pulsar-client-go/pulsar"
)

var (
	client pulsar.Client
	once   sync.Once
)

// NewPulsarDS ...
func NewPulsarDS(dsURL string) (pulsar.Client, error) {
	once.Do(func() {
		client, err = pulsar.NewClient(pulsar.ClientOptions{
			URL: dsURL,
		})
	})
	if err != nil {
		return nil, err
	}

	return client, err
}

// ClosePulsarDS ...
func ClosePulsarDS() {
	client.Close()
}
