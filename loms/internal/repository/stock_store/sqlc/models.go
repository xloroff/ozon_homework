// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package sqlc

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type OrderStatusType string

const (
	OrderStatusTypeNew             OrderStatusType = "new"
	OrderStatusTypeAwaitingPayment OrderStatusType = "awaiting_payment"
	OrderStatusTypePayed           OrderStatusType = "payed"
	OrderStatusTypeCancelled       OrderStatusType = "cancelled"
	OrderStatusTypeFailed          OrderStatusType = "failed"
)

func (e *OrderStatusType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = OrderStatusType(s)
	case string:
		*e = OrderStatusType(s)
	default:
		return fmt.Errorf("unsupported scan type for OrderStatusType: %T", src)
	}
	return nil
}

type NullOrderStatusType struct {
	OrderStatusType OrderStatusType
	Valid           bool // Valid is true if OrderStatusType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullOrderStatusType) Scan(value interface{}) error {
	if value == nil {
		ns.OrderStatusType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.OrderStatusType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullOrderStatusType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.OrderStatusType), nil
}

type Order struct {
	ID        int64
	User      int64
	Status    OrderStatusType
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type OrderItem struct {
	ID      int64
	OrderID int64
	Sku     int64
	Count   int32
}

type Stock struct {
	Sku        int64
	TotalCount int32
	Reserved   int32
}
