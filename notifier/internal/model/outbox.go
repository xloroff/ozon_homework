package model

import (
	"time"
)

// OrderEventPayload полезная нагрузка события по заказу которая будет получена из очереди.
type OrderEventPayload struct {
	ID       string    `json:"id"`
	Time     time.Time `json:"time"`
	EntityID string    `json:"entity_id"`
	Payload  Order     `json:"payload"`
}
