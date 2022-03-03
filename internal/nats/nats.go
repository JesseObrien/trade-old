package nats

import (
	nats "github.com/nats-io/nats.go"
)

// NewJSONConnection gives you back a json encoded nats connection
func NewJSONConnection(natsURL string) (*nats.EncodedConn, error) {
	nc, err := nats.Connect(natsURL)

	if err != nil {
		return nil, err
	}

	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	return c, err
}
