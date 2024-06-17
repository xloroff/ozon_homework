package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Querier интерфейс для выполнения запросов к БД.
type Querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) error
	BeginFuncWithTx(ctx context.Context, f func(pgx.Tx) error) error
	Ping(ctx context.Context) error
	Close() error
}

// Exec выполняет запрос к БД при этом воззвращая соединение в пул.
func (p *pool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	res, err := p.pool.Exec(ctx, sql, args...)
	if err != nil {
		return res, fmt.Errorf("Ошибка выполнения запроса - %w", err)
	}

	return res, nil
}

// Query выполняет запрос к БД.
func (p *pool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	res, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return res, fmt.Errorf("Ошибка выполнения запроса - %w", err)
	}

	return res, nil
}

// QueryRow выполняет запрос к БД, но возвращает только первую строку.
func (p *pool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

// BeginFunc выполняет запрос к БД если вполнение было без ошибок, вызовет commit.
func (p *pool) BeginFunc(ctx context.Context, f func(pgx.Tx) error) error {
	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("Ошибка возвращения соединения пула - %w", err)
	}
	defer conn.Release()

	err = pgx.BeginFunc(ctx, conn, f)
	if err != nil {
		return fmt.Errorf("Ошибка выполнения транзакции/запроса - %w", err)
	}

	return nil
}

// BeginFuncWithTx выполняет запрос к БД если вполнение было без ошибок, вызовет commit.
func (p *pool) BeginFuncWithTx(ctx context.Context, f func(pgx.Tx) error) error {
	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("Ошибка возвращения соединения пула - %w", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("Ошибка создания транзакции - %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	err = f(tx)
	if err != nil {
		return fmt.Errorf("Ошибка выполнения транзакции/запроса - %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("Ошибка фиксации транзакции - %w", err)
	}

	return nil
}

// Ping проверяет соединение с БД.
func (p *pool) Ping(ctx context.Context) error {
	err := p.pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("Ошибка связи - %w", err)
	}

	return nil
}

// Close закрывает соединение с БД.
func (p *pool) Close() error {
	p.pool.Close()
	return nil
}
