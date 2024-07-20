package model

import "errors"

var (
	// ErrOrderNotFound ошибка заказа.
	ErrOrderNotFound = errors.New("Заказ не найден")

	// ErrOrderCancel ошибка отмены заказа в текущем статусе.
	ErrOrderCancel = errors.New("Нельзя отменить заказ в этом статусе")

	// ErrReserveNotFound ошибка поиска остатков для товара.
	ErrReserveNotFound = errors.New("Остатки по товару не найдены")
)
