package restycli

import "errors"

var (
	// ErrNotFound ошибка отсутствия товара на сервисе продуктов.
	ErrNotFound = errors.New("Товар не найден")
	// ErrCartEmpty корзина пользователя пуста.
	ErrCartEmpty = errors.New("Корзина пользователя пуста")
	// ErrTooManyRequests ошибка при ответе сервиса продуктов что у него слишком много запросов.
	ErrTooManyRequests = errors.New("Cлишком много запросов")
	// ErrUnknownError неизвестная ошибка.
	ErrUnknownError = errors.New("Неизвестная ошибка")
)