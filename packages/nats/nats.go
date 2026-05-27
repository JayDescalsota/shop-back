package natss

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func NewConnection(url string) (*nats.Conn, jetstream.JetStream, error) {
	nc, err := nats.Connect(url,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(2*time.Second),
		nats.Timeout(5*time.Second),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("nats connect: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, nil, fmt.Errorf("jetstream new: %w", err)
	}

	log.Println("connected to nats jetstream")
	return nc, js, nil
}

func EnsureStream(js jetstream.JetStream, streamName string, subjects ...string) error {
	_, err := js.CreateStream(nil, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: subjects,
		MaxAge:   7 * 24 * time.Hour,
		Storage:  jetstream.FileStorage,
		Replicas: 1,
	})
	if err != nil {
		return fmt.Errorf("create stream %s: %w", streamName, err)
	}
	return nil
}
