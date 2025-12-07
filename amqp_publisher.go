package main

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/ride4Low/contracts/events"
)

type MessagePublisher interface {
	PublishMessage(ctx context.Context, routingKey string, message events.AmqpMessage) error
}

type EventPublisher interface {
	PublishPaymentSuccess(ctx context.Context, event *events.PaymentStatusUpdateData) error
}

type AmqpPublisher struct {
	publisher MessagePublisher
}

func NewAmqpPublisher(publisher MessagePublisher) *AmqpPublisher {
	return &AmqpPublisher{publisher: publisher}
}

func (p *AmqpPublisher) PublishPaymentSuccess(ctx context.Context, event *events.PaymentStatusUpdateData) error {
	payloadBytes, err := sonic.Marshal(event)
	if err != nil {
		return err
	}

	amqpMessage := events.AmqpMessage{
		OwnerID: event.UserID,
		Data:    payloadBytes,
	}

	return p.publisher.PublishMessage(ctx, events.PaymentEventSuccess, amqpMessage)
}
