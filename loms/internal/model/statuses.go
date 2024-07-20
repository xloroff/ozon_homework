package model

const (
	// OrderStatusNew статус нового заказа.
	OrderStatusNew = "new"

	// OrderStatusAwaitingPayment статус заказа ожидающего оплаты.
	OrderStatusAwaitingPayment = "awaiting_payment"

	// OrderStatusPayed заказ оплачен.
	OrderStatusPayed = "payed"

	// OrderStatusCancelled заказ отменен.
	OrderStatusCancelled = "cancelled"

	// OrderStatusFailed заказ помечен как неудачный.
	OrderStatusFailed = "failed"
)
