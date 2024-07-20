package model

// Order весь заказ.
type Order struct {
	ID     int64
	User   int64
	Status string
	Items  OrderItems `json:",omitempty"`
}

// OrderItem отдельный итем в заказе.
type OrderItem struct {
	ID    int64
	Sku   int64
	Count uint16
}

// OrderItems лист итемов заказа.
type OrderItems []*OrderItem

// AllOrderItems все хранилише заказов в памяти.
type AllOrderItems map[int64]*Order
