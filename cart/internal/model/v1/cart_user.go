package v1

const (
	// UsrID упорядоченный нейминг в функциях идентификатора пользователя.
	UsrID = "user_id"
	// SkuID упорядоченный нейминг в функциях идентификатора товара.
	SkuID = "sku_id"
)

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

// ApiResponse в ТЗ не указано, но буду в ответ ошибки прокидывать если не будет ошибки - будет пустой ответ.
type ApiResponse struct {
	Error string `json:"error,omitempty"`
}

// UserIdintyfier идентификатор пользователя в сервисе.
type UserIdintyfier struct {
	UserID int64 `validate:"required|int|min:1" message:"Поле:{user_id} обязательно для указания в строке запроса - минимальное значение 1"`
}

// UsrSkuID идентификатор товара.
type UsrSkuID struct {
	SkuID int64 `validate:"required|int|min:1" message:"Поле:{sku_id} обязательно для указания в строке запроса - минимальное значение 1"`
}