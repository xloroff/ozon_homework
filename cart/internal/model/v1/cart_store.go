package v1

// CartItem для получения числа итемов.
type CartItem struct {
	Count uint16
}

// CartItems для хранения в памяти числа покупок корзин.
type CartItems map[int64]*CartItem

type Cart struct {
	Items CartItems
}

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