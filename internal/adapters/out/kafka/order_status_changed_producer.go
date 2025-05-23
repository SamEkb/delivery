package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/delivery/internal/core/domain/model/order"
	"github.com/delivery/internal/core/ports"
	"github.com/delivery/internal/generated/events/queues/orderstatuschangedpb"
	"github.com/delivery/internal/pkg/ddd"
	"github.com/delivery/internal/pkg/errs"
	"google.golang.org/protobuf/proto"
)

type orderStatusChangedProducer struct {
	topic    string
	producer sarama.SyncProducer
}

func NewOrderStatusChangedProducer(brokers []string, topic string) (ports.OrderProducer, error) {
	if brokers == nil || len(brokers) == 0 {
		return nil, errs.NewValueIsRequiredError("brokers")
	}
	if topic == "" {
		return nil, errs.NewValueIsRequiredError("topic")
	}

	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V3_4_0_0
	saramaCfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama sync producer: %w", err)
	}

	return &orderStatusChangedProducer{
		topic:    topic,
		producer: producer,
	}, nil
}

func (o *orderStatusChangedProducer) Publish(ctx context.Context, domainEvent ddd.DomainEvent) error {
	if domainEvent == nil {
		return errs.NewValueIsRequiredError("event")
	}

	completedDomainEvent, ok := domainEvent.(*order.StatusChangedDomainEvent)
	if !ok {
		return fmt.Errorf("invalid event type: %T, expected: %T", domainEvent, &order.StatusChangedDomainEvent{})
	}

	integrationEvent, err := mapDomainEventToIntegrationEvent(completedDomainEvent)
	if err != nil {
		return fmt.Errorf("failed to map domain event to integration event: %w", err)
	}

	eventBytes, err := proto.Marshal(integrationEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal integration event: %w", err)
	}

	message := &sarama.ProducerMessage{
		Topic: o.topic,
		Key:   sarama.StringEncoder(completedDomainEvent.OrderID.String()),
		Value: sarama.ByteEncoder(eventBytes),
	}

	resultCh := make(chan error, 1)

	go func() {
		partition, offset, errSend := o.producer.SendMessage(message)
		if errSend == nil {
			log.Printf("Message for order %s sent successfully to topic %s, partition %d, offset %d",
				completedDomainEvent.OrderID.String(), o.topic, partition, offset)
		}
		resultCh <- errSend
		close(resultCh)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case errSend := <-resultCh:
		if errSend != nil {
			return fmt.Errorf("failed to send message to kafka: %w", errSend)
		}
		return nil
	}
}

func (o *orderStatusChangedProducer) Close() error {
	if err := o.producer.Close(); err != nil {
		return fmt.Errorf("error closing sarama producer: %w", err)
	}
	return nil
}

func mapDomainEventToIntegrationEvent(event *order.StatusChangedDomainEvent) (*orderstatuschangedpb.OrderStatusChangedIntegrationEvent, error) {
	if event == nil {
		return nil, errs.NewValueIsRequiredError("event (StatusChangedDomainEvent)")
	}

	status, ok := orderstatuschangedpb.OrderStatus_value[event.OrderStatus.String()]
	if !ok {
		return nil, fmt.Errorf("invalid order status: %s", event.OrderStatus.String())
	}

	integrationEvent := orderstatuschangedpb.OrderStatusChangedIntegrationEvent{
		OrderId:     event.OrderID.String(),
		OrderStatus: orderstatuschangedpb.OrderStatus(status),
	}

	return &integrationEvent, nil
}
