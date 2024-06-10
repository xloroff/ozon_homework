package model

// FullUserCart корзина пользователя.
type FullUserCart struct {
	Items      []*UserCartItem `json:"items"`
	TotalPrice uint32          `json:"total_price"`
}

// UserCartItem одна позиция в корзине.
type UserCartItem struct {
	SkuID int64  `json:"sku_id"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}
