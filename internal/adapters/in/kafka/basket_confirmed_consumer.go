package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/delivery/internal/core/application/usecases/commands"
	"github.com/delivery/internal/generated/messages/queues/basketconfirmedpb"
	"github.com/google/uuid"
)

type BasketConfirmedConsumer interface {
	Consume() error
	Close() error
}

var _ BasketConfirmedConsumer = &basketConfirmedConsumer{}
var _ sarama.ConsumerGroupHandler = &basketConfirmedConsumer{}

type basketConfirmedConsumer struct {
	topic                     string
	consumerGroup             sarama.ConsumerGroup
	createOrderCommandHandler commands.CreateOrderHandler
	ctx                       context.Context
	cancel                    context.CancelFunc
}

func NewConsumer(
	brokers []string,
	topic string,
	group string,
	createOrderCommandHandler commands.CreateOrderHandler,
) (BasketConfirmedConsumer, error) {

	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V3_4_0_0
	saramaCfg.Consumer.Return.Errors = true
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &basketConfirmedConsumer{
		topic:                     topic,
		consumerGroup:             consumerGroup,
		createOrderCommandHandler: createOrderCommandHandler,
		ctx:                       ctx,
		cancel:                    cancel,
	}, nil
}

func (b *basketConfirmedConsumer) Consume() error {
	for {
		if err := b.consumerGroup.Consume(b.ctx, []string{b.topic}, b); err != nil {
			return err
		}
		if b.ctx.Err() != nil {
			return nil
		}
	}
}

func (b *basketConfirmedConsumer) Close() error {
	b.cancel()
	return b.consumerGroup.Close()
}

func (b *basketConfirmedConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (b *basketConfirmedConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (b *basketConfirmedConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("Messages channel closed for claim topic: %s, partition: %d", claim.Topic(), claim.Partition())
				return nil
			}

			var event basketconfirmedpb.BasketConfirmedIntegrationEvent
			if err := json.Unmarshal(message.Value, &event); err != nil {
				log.Printf("Failed to unmarshal message for topic %s, partition %d, offset %d: %v. Skipping message.",
					message.Topic, message.Partition, message.Offset, err)
				session.MarkMessage(message, "")
				continue
			}

			parsedBasketID, err := uuid.Parse(event.BasketId)
			if err != nil {
				log.Printf("Failed to parse BasketId '%s' as UUID for topic %s, partition %d, offset %d: %v. Skipping message.",
					event.BasketId, message.Topic, message.Partition, message.Offset, err)
				session.MarkMessage(message, "")
				continue
			}

			command, err := commands.NewCreateOrderCommand(
				parsedBasketID,
				event.GetAddress().GetStreet(),
				int(event.GetVolume()),
			)

			if err != nil {
				log.Printf("Failed to create NewCreateOrderCommand for topic %s, partition %d, offset %d: %v. Skipping message.",
					message.Topic, message.Partition, message.Offset, err)
				session.MarkMessage(message, "")
				continue
			}

			handlerCtx := session.Context()
			if err := b.createOrderCommandHandler.Handle(handlerCtx, command); err != nil {
				log.Printf("Failed to handle CreateOrderCommand for topic %s, partition %d, offset %d: %v. Message will be reprocessed.",
					message.Topic, message.Partition, message.Offset, err)
				return fmt.Errorf("failed to process message offset %d: %w", message.Offset, err)
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			log.Printf("Session context done for claim topic: %s, partition: %d. Exiting ConsumeClaim.", claim.Topic(), claim.Partition())
			return nil
		}
	}
}
