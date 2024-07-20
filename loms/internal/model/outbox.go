package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Outbox сообщение хранимое в БД для кафки.
type Outbox struct {
	ID        string
	CreatedAt sql.NullTime
	EntityID  string
	Status    string
	Payload   string
	Metadata  json.RawMessage
}

// OrderEventPayload полезная нагрузка события по заказу которая будет отправляться в очередь.
type OrderEventPayload struct {
	ID       string    `json:"id"`
	Time     time.Time `json:"time"`
	EntityID string    `json:"entity_id"`
	Payload  Order     `json:"payload"`
}

// Metadata метаданные о сообщении для хранении в БД outbox.
type Metadata struct {
	TraceID string `json:"trace_id"`
	SpanID  string `json:"span_id"`
}
