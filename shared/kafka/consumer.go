package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/SureshAmal/NimbusU-backend/shared/config"
	"github.com/SureshAmal/NimbusU-backend/shared/logger"
	"go.uber.org/zap"
)

// MessageHandler is a function that processes a Kafka message
type MessageHandler func(ctx context.Context, message []byte) error

// Consumer wraps Sarama consumer group
type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	handler       MessageHandler
	topics        []string
}

// consumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	handler MessageHandler
}

func (h consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		logger.Debug("Message consumed",
			zap.String("topic", message.Topic),
			zap.Int32("partition", message.Partition),
			zap.Int64("offset", message.Offset),
		)

		// Process message
		err := h.handler(session.Context(), message.Value)
		if err != nil {
			logger.Error("Error processing message",
				zap.String("topic", message.Topic),
				zap.Int64("offset", message.Offset),
				zap.Error(err),
			)
			// Continue processing other messages even if one fails
			continue
		}

		// Mark message as processed
		session.MarkMessage(message, "")
	}

	return nil
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg config.KafkaConfig, topics []string, handler MessageHandler) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.ConsumerGroup, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer group: %w", err)
	}

	logger.Info("Kafka consumer created",
		zap.Strings("brokers", cfg.Brokers),
		zap.String("group", cfg.ConsumerGroup),
		zap.Strings("topics", topics),
	)

	return &Consumer{
		consumerGroup: consumerGroup,
		handler:       handler,
		topics:        topics,
	}, nil
}

// Start starts consuming messages from Kafka
func (c *Consumer) Start(ctx context.Context) error {
	handler := consumerGroupHandler{handler: c.handler}

	for {
		// Check if context is cancelled
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Consume messages
		err := c.consumerGroup.Consume(ctx, c.topics, handler)
		if err != nil {
			logger.Error("Error consuming messages", zap.Error(err))
			return err
		}

		// Check if context was cancelled during consume
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

// Close closes the Kafka consumer
func (c *Consumer) Close() error {
	if c.consumerGroup != nil {
		err := c.consumerGroup.Close()
		if err != nil {
			logger.Error("Error closing Kafka consumer", zap.Error(err))
			return err
		}
		logger.Info("Kafka consumer closed")
	}
	return nil
}

// UnmarshalEvent unmarshals a JSON event from bytes
func UnmarshalEvent(data []byte, event interface{}) error {
	return json.Unmarshal(data, event)
}
