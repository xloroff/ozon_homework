package kafkaproducer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
)

// Producer методы доступные для работы продюссера с Kafka.
type Producer interface {
	Send(ctx context.Context, key string, message any) error
	Close() error
}

type producer struct {
	logger   logger.Logger
	producer sarama.SyncProducer
	topic    string
}

// NewProducer создает новый экземпляр Producer.
func NewProducer(addr, topic string, l logger.Logger, opts ...ProducerOption) (Producer, error) {
	cfg := NewConfig(opts...)

	syncProducer, err := sarama.NewSyncProducer([]string{addr}, cfg)
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания SyncProducer - %w", err)
	}

	return &producer{
		logger:   l,
		producer: syncProducer,
		topic:    topic,
	}, nil
}

// SendMessage отправляет сообщение в топик.
func (p *producer) Send(ctx context.Context, key string, message any) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("Ошибка форматирования сообщения - %w", err)
	}

	pm := &sarama.ProducerMessage{
		Headers: []sarama.RecordHeader{
			{
				Key:   sarama.ByteEncoder("app-name"),
				Value: sarama.ByteEncoder(config.AppName),
			},
			{
				Key:   sarama.ByteEncoder("x-trace-id"),
				Value: sarama.ByteEncoder(tracer.GetTraceID(ctx)),
			},
			{
				Key:   sarama.ByteEncoder("x-span-id"),
				Value: sarama.ByteEncoder(tracer.GetSpanID(ctx)),
			},
		},
		Timestamp: time.Now(),
		Topic:     p.topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.StringEncoder(msg),
	}

	_, _, err = p.producer.SendMessage(pm)
	if err != nil {
		return fmt.Errorf("Ошибка отправки сообщения в топик - %s: %w", p.topic, err)
	}

	return nil
}

// Close закрывает producer.
func (p *producer) Close() error {
	ctx := context.Background()

	err := p.producer.Close()
	if err != nil {
		return fmt.Errorf("Ошибка закрытия producer - %w", err)
	}

	p.logger.Info(ctx, "kafka producer is closed successfully")

	return nil
}
