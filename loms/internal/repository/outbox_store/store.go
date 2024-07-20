package outboxstore

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/outbox_store/sqlc"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/pkg/db"
)

const (
	repName = "OutboxStore"
)

// Storage имплементирует методы управления хранилищем памяти.
type Storage interface {
	GetEvents(ctx context.Context) (msgs []*sqlc.OutboxRow, err error)
	SetStatus(ctx context.Context, msg *sqlc.OutboxRow) error
	AddMessage(ctx context.Context, tx pgx.Tx, message *model.Outbox) error
}

type outboxStorage struct {
	ctx      context.Context
	lockTime int
	data     db.ClientBD
	logger   logger.Logger
}

// NewOutboxStorage создает хранилище/клиента кафки в БД.
func NewOutboxStorage(ctx context.Context, l logger.Logger, bdCli db.ClientBD, lockTime int) (Storage, error) {
	return &outboxStorage{
		ctx:      ctx,
		lockTime: lockTime,
		data:     bdCli,
		logger:   l,
	}, nil
}

func toPGText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{}
	}

	return pgtype.Text{String: s, Valid: true}
}
