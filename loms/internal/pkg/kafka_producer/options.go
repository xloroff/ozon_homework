package kafkaproducer

import (
	"time"

	"github.com/IBM/sarama"
)

const (
	defaultMaxRetry   = 30
	defaultBackoff    = 5
	defaultMaxOpenReq = 1
)

// ProducerOption кастомный тип опций конфигурации для продюсера.
type ProducerOption func(c *sarama.Config)

// NewConfig конфигурация для продюсера с дефолтными параметрами.
func NewConfig(opts ...ProducerOption) *sarama.Config {
	c := sarama.NewConfig()

	// Алгоритм выбора партиции.
	c.Producer.Partitioner = sarama.NewHashPartitioner
	// Acsk параметр.
	c.Producer.RequiredAcks = sarama.WaitForAll
	// Cемантика exactly once.
	c.Producer.Idempotent = false
	// Число попыток отправить сообщение.
	c.Producer.Retry.Max = defaultMaxRetry
	// Интервалы между попытками отправить сообщение.
	c.Producer.Retry.Backoff = defaultBackoff * time.Millisecond
	// Количество соединений которое может быть одновремнно открыто. Одно соединение гарантирует порядорк доставки сообщений.
	c.Net.MaxOpenRequests = defaultMaxOpenReq
	// Используемое сжатие.
	c.Producer.Compression = sarama.CompressionGZIP
	// Уровень сжатия.
	c.Producer.CompressionLevel = sarama.CompressionLevelDefault
	// Использование конфигурации для создания SyncProducer. Устанавливаем, оба в true, гарантируя, что читать не будем и отдаем под капот sarama всё.
	c.Producer.Return.Successes = true
	c.Producer.Return.Errors = true
	// Реконнекты к кафке если она отвалилась
	c.Metadata.Retry.Max = defaultMaxRetry
	c.Metadata.Retry.Backoff = defaultBackoff * time.Millisecond

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithRequiredAcks переопределяет опцию гарантии доставки сообщений.
func WithRequiredAcks(acks sarama.RequiredAcks) ProducerOption {
	return func(c *sarama.Config) {
		c.Producer.RequiredAcks = acks
	}
}

// WithIdempotent переопределяет опцию - записана ровно одна копия сообщения.
func WithIdempotent() ProducerOption {
	return func(c *sarama.Config) {
		c.Producer.Idempotent = true
	}
}

// WithRetryMax переопределяет опцию максимальное число попыток отправить сообщение.
func WithRetryMax(n int) ProducerOption {
	return func(c *sarama.Config) {
		c.Producer.Retry.Max = n
	}
}

// WithRetryBackoff переопределяет интервал между попытками отправить сообщение.
func WithRetryBackoff(n int) ProducerOption {
	return func(c *sarama.Config) {
		c.Producer.Retry.Backoff = time.Duration(n) * time.Millisecond
	}
}

// WithMaxOpenRequests переопределяет максимальное число открытых соединений.
func WithMaxOpenRequests(n int) ProducerOption {
	return func(c *sarama.Config) {
		c.Net.MaxOpenRequests = n
	}
}
