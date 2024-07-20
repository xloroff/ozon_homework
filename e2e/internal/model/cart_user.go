package model

// AddItem собираем данные всего входящего запроса (все параметры) на добавление итема в корзину.
type AddItem struct {
	UserIdintyfier
	UsrSkuID
	AddItemBody
}

// DelItem собираем данные входящего запроса (все параметры) на удаление итема из корзины.
type DelItem struct {
	UserIdintyfier
	UsrSkuID
}

// AddItemBody запрос на добавление товара в корзину.
type AddItemBody struct {
	Count uint16 `json:"count,omitempty" validate:"required|int|min:1" message:"Поле:{count} обязательно для указания в теле запроса - минимальное значение 1"`
}

// UserIdintyfier идентификатор пользователя в сервисе.
type UserIdintyfier struct {
	UserID int64 `validate:"required|int|min:1" message:"Поле:{user_id} обязательно для указания в строке запроса - минимальное значение 1"`
}

// UsrSkuID идентификатор товара.
type UsrSkuID struct {
	SkuID int64 `validate:"required|int|min:1" message:"Поле:{sku_id} обязательно для указания в строке запроса - минимальное значение 1"`
}
