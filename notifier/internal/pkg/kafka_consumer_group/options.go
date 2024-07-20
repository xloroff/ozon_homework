package kafkaconsumergroup

import (
	"time"

	"github.com/IBM/sarama"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/config"
)

const (
	defaultHeartbeatInterval         = 3
	defaultSessionTimeout            = 60
	defaultRebalanceTimeout          = 60
	defaultOffsetsAutoCommitInterval = 5
)

// ConsumerGroupOption кастомный тип опций конфигурации для консьюмер группы.
type ConsumerGroupOption func(c *sarama.Config)

// NewConfig конфигурация для консьюмер группы с дефолтными параметрами.
func NewConfig(opts ...ConsumerGroupOption) *sarama.Config {
	saramaConfig := sarama.NewConfig()

	saramaConfig.ClientID = config.AppName

	saramaConfig.Version = sarama.MaxVersion
	// Читаем сообщения с самого начала.
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	// Пропуск невалидных сдвигов, в случае убегания offset.
	saramaConfig.Consumer.Group.ResetInvalidOffsets = true
	// Интервал отправки признаков жизни.
	saramaConfig.Consumer.Group.Heartbeat.Interval = defaultHeartbeatInterval * time.Second
	// Таймаут сессии с кафкой.
	saramaConfig.Consumer.Group.Session.Timeout = defaultSessionTimeout * time.Second
	// Таймаут ребалансировки.
	saramaConfig.Consumer.Group.Rebalance.Timeout = defaultRebalanceTimeout * time.Second
	// Будем обрабатывать ошибки.
	saramaConfig.Consumer.Return.Errors = true
	// Отключаем автокоммит оффсета, но оставляем интервал на случай если захотим включить.
	saramaConfig.Consumer.Offsets.AutoCommit.Enable = false
	saramaConfig.Consumer.Offsets.AutoCommit.Interval = defaultOffsetsAutoCommitInterval * time.Second

	for _, opt := range opts {
		opt(saramaConfig)
	}

	return saramaConfig
}

// WithOffsetsAutoCommitInterval переопределяет интервал автокоммита.
func WithOffsetsAutoCommitInterval(n int) ConsumerGroupOption {
	return func(c *sarama.Config) {
		c.Consumer.Offsets.AutoCommit.Interval = time.Duration(n) * time.Second
	}
}

// WithHeartbeatInterval переопределяет интервал heartbeat.
func WithHeartbeatInterval(n int) ConsumerGroupOption {
	return func(c *sarama.Config) {
		c.Consumer.Group.Heartbeat.Interval = time.Duration(n) * time.Second
	}
}

// WithSessionTimeout переопределяет таймаут сессии.
func WithSessionTimeout(n int) ConsumerGroupOption {
	return func(c *sarama.Config) {
		c.Consumer.Group.Session.Timeout = time.Duration(n) * time.Second
	}
}

// WithRebalanceTimeout переопределяет таймаут ребалансировки.
func WithRebalanceTimeout(n int) ConsumerGroupOption {
	return func(c *sarama.Config) {
		c.Consumer.Group.Rebalance.Timeout = time.Duration(n) * time.Second
	}
}

// WithAutoCommit переопределяет настройку автокоммита.
func WithAutoCommit(enable bool) ConsumerGroupOption {
	return func(c *sarama.Config) {
		c.Consumer.Offsets.AutoCommit.Enable = enable
	}
}
