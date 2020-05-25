package mb

import (
	"sync"

	"github.com/apache/pulsar-client-go/pulsar"
)

var (
	client pulsar.Client
	once   sync.Once
)

// NewPulsarMB ...
func NewPulsarMB(mbURL string) (pulsar.Client, error) {
	once.Do(func() {
		client, err = pulsar.NewClient(pulsar.ClientOptions{
			URL: mbURL,
		})
	})
	if err != nil {
		return nil, err
	}

	return client, err
}

// ClosePulsarMB ...
func ClosePulsarMB() {
	client.Close()
}
