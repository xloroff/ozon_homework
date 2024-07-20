package kafkaconsumergroup

import (
	"context"
	"fmt"
	"strings"

	"github.com/IBM/sarama"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/consumers/order"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/pkg/logger"
)

// ConsumerGroup методы доступные для управления группой консьюмеров.
type ConsumerGroup interface {
	Start() error
	Stop() error
}

type consumerGroup struct {
	ctx       context.Context
	logger    logger.Logger
	consGroup sarama.ConsumerGroup
	consHadl  sarama.ConsumerGroupHandler
	topicList []string
}

// NewConsumerGroup создает группу консьюмеров kafka.
func NewConsumerGroup(ctx context.Context, l logger.Logger, cnfg *config.KafkaSettings, consumer orderconsumer.Consumer, opts ...ConsumerGroupOption) (ConsumerGroup, error) {
	consumerConfig := NewConfig(opts...)

	cg, err := sarama.NewConsumerGroup([]string{cnfg.KafkaAddress}, cnfg.KafkaConsumerGroupName, consumerConfig)
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания группы консьюмеров - %w", err)
	}

	return &consumerGroup{
		ctx:       ctx,
		logger:    l,
		consGroup: cg,
		consHadl:  newConsumerGroupHandler(ctx, l, consumer.OrderHandle),
		topicList: []string{cnfg.KafkaTopic},
	}, nil
}

// Start запускает группу консьюмеров.
func (cg *consumerGroup) Start() error {
	for {
		if err := cg.consGroup.Consume(cg.ctx, cg.topicList, cg.consHadl); err != nil {
			if err != sarama.ErrClosedConsumerGroup {
				cg.logger.Errorf(cg.ctx, "KafkaConsumerGroup.Start: Ошибка наблюдени за топиком/топиками - %v: %v", strings.Join(cg.topicList, ", "), err)
				return fmt.Errorf("Ошибка наблюдени за топиком/топиками - %v: %w", strings.Join(cg.topicList, ", "), err)
			}
		}
	}
}

// Stop останавливает группу консьюмеров.
func (cg *consumerGroup) Stop() error {
	err := cg.consGroup.Close()
	if err != nil {
		cg.logger.Errorf(cg.ctx, "KafkaConsumerGroup.Stop: Ошибка остановки группы констюмеров - %v", err)
		return fmt.Errorf("Ошибка остановки группы констюмеров - %w", err)
	}

	cg.logger.Info(cg.ctx, "Остановка группы констюмеров завершена...")

	return nil
}
