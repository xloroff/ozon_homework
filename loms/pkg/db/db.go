package db

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/jackc/pgx/v5/pgxpool"
	// Добавляем pgx драйвер для миграции как подмену stdlib.
	_ "github.com/jackc/pgx/v5/stdlib"
)

// ClientBD репрезентационный клиент для работы с БД.
type ClientBD interface {
	GetReaderPool() Querier
	GetWriterPool() Querier
	GetMasterPool() Querier
	GetSyncPool() Querier
	Close() error
}

type client struct {
	ctx               context.Context
	mPool             Querier
	sPool             Querier
	readerPoolCounter atomic.Uint64
}

type pool struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

// NewClient создаёт новый экземпляр клиента с подключением к нескольким пулам.
func NewClient(ctx context.Context, masterConnStr, syncConnStr string) (ClientBD, error) {
	masterPool, err := NewPool(ctx, masterConnStr)
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания master pool - %w", err)
	}

	syncPool, err := NewPool(ctx, syncConnStr)
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания sync pool - %w", err)
	}

	return &client{
		ctx:   ctx,
		mPool: masterPool,
		sPool: syncPool,
	}, nil
}

// NewPool создает пул/ресурс для подключения к БД.
func NewPool(ctx context.Context, cnStr string) (Querier, error) {
	p, err := pgxpool.New(ctx, cnStr)
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания пула, строка подключения - %s: %w", cnStr, err)
	}

	err = p.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("Ошибка проверки связи с БД, строка подключения -  %s:  %w", cnStr, err)
	}

	return &pool{
		ctx:  ctx,
		pool: p,
	}, nil
}

// GetReaderPool возвращает пул для чтения данных из БД.
func (c *client) GetReaderPool() Querier {
	res := c.readerPoolCounter.Add(1)
	if res%2 == 0 {
		return c.mPool
	}

	return c.sPool
}

// GetWriterPool возвращает пул для записи данных в БД.
func (c *client) GetWriterPool() Querier {
	return c.GetMasterPool()
}

// GetMasterPool возвращает мастер пул для чтения данных из БД.
func (c *client) GetMasterPool() Querier {
	return c.mPool
}

// GetSyncPool возвращает синк-пул для чтения данных из БД.
func (c *client) GetSyncPool() Querier {
	return c.sPool
}

// Close освобождает пулы.
func (c *client) Close() error {
	err := c.mPool.Close()
	if err != nil {
		return fmt.Errorf("Ошибка закрытия master pool - %w", err)
	}

	err = c.sPool.Close()
	if err != nil {
		return fmt.Errorf("Ошибка закрытия sync pool - %w", err)
	}

	return nil
}
