package outboxservice

import (
	"context"
	"sync"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/kafka_producer"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/outbox_store"
)

// Service управление сервисом по работе с брокером.
type Service interface {
	Start(ctx context.Context)
	Stop() error
}

type service struct {
	ctx            context.Context
	periodSend     int
	outboxstore    outboxstore.Storage
	outboxProducer kafkaproducer.Producer
	logger         logger.Logger
	sendStatuses   chan struct{}
	wg             sync.WaitGroup
}

// NewService создает экземпляр сервиса по работе с брокером сообщений.
func NewService(ctx context.Context, l logger.Logger, period int, s outboxstore.Storage, p kafkaproducer.Producer) Service {
	return &service{
		ctx:            ctx,
		periodSend:     period,
		outboxstore:    s,
		outboxProducer: p,
		logger:         l,
		sendStatuses:   make(chan struct{}),
	}
}
