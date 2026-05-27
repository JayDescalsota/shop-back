package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Envelope struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Source        string                 `json:"source"`
	TenantID      string                 `json:"tenant_id"`
	UserID        string                 `json:"user_id,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	CorrelationID string                 `json:"correlation_id"`
	Payload       map[string]interface{} `json:"payload"`
}

type Publisher struct {
	js     jetstream.JetStream
	source string
}

func NewPublisher(js jetstream.JetStream, source string) *Publisher {
	return &Publisher{js: js, source: source}
}

func (p *Publisher) Publish(eventType string, tenantID string, payload map[string]interface{}) (string, error) {
	evt := Envelope{
		ID:        uuid.New().String(),
		Type:      eventType,
		Source:    p.source,
		TenantID:  tenantID,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}

	data, err := json.Marshal(evt)
	if err != nil {
		return "", err
	}

	_, err = p.js.Publish(nil, eventType, data)
	if err != nil {
		return "", err
	}

	return evt.ID, nil
}

func (p *Publisher) PublishWithUser(eventType, tenantID, userID string, payload map[string]interface{}) (string, error) {
	evt := Envelope{
		ID:        uuid.New().String(),
		Type:      eventType,
		Source:    p.source,
		TenantID:  tenantID,
		UserID:    userID,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}

	data, err := json.Marshal(evt)
	if err != nil {
		return "", err
	}

	_, err = p.js.Publish(nil, eventType, data)
	if err != nil {
		return "", err
	}

	return evt.ID, nil
}

func (p *Publisher) PublishWithCorrelation(eventType, tenantID, correlationID string, payload map[string]interface{}) (string, error) {
	evt := Envelope{
		ID:            uuid.New().String(),
		Type:          eventType,
		Source:        p.source,
		TenantID:      tenantID,
		CorrelationID: correlationID,
		Timestamp:     time.Now().UTC(),
		Payload:       payload,
	}

	data, err := json.Marshal(evt)
	if err != nil {
		return "", err
	}

	_, err = p.js.Publish(nil, eventType, data)
	if err != nil {
		return "", err
	}

	return evt.ID, nil
}

type Handler func(Envelope) error

func Subscribe(js jetstream.JetStream, subject string, handler Handler) error {
	cons, err := js.CreateOrUpdateConsumer(nil, jetstream.ConsumerConfig{
		Name:          subject + "-consumer",
		FilterSubject: subject,
		AckPolicy:    jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return err
	}

	cc := make(chan struct{})
	_ = nats.NewInbox()

	go func() {
		for {
			msg, err := cons.Next()
			if err != nil {
				continue
			}

			var evt Envelope
			if err := json.Unmarshal(msg.Data(), &evt); err != nil {
				msg.Nak()
				continue
			}

			if err := handler(evt); err != nil {
				msg.Nak()
				continue
			}

			msg.Ack()
		}
	}()

	<-cc
	return nil
}
