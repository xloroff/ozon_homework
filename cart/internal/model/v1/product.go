package v1

// ProductReq запрос сервису продуктов.
type ProductReq struct {
	Token string `json:"token"`
	Sku   int64  `json:"sku"`
}

// ProductResp ответ от сервиса продуктов.
type ProductResp struct {
	Name  string `json:"name" validate:"required|min_len:1" message:"Имя продукта {name} полученное от сервиса продуктов некорректно"`
	Price uint32 `json:"price" validate:"required|int|min:0" message:"Стоимость продукта {price} полученное от сервиса продуктов некорректно"`
}