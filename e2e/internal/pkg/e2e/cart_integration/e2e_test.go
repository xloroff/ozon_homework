//go:build e2e

package cartintegration

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/e2e/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/e2e/internal/pkg/client/cart_cli"
)

type Suite struct {
	suite.Suite
	ctx  context.Context
	cart cartcli.CartClient
}

func TestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &Suite{})
}

func (s *Suite) SetupSuite() {
	s.ctx = context.Background()

	s.cart = cartcli.NewCartClient(fmt.Sprintf("http://%s", os.Getenv("CART_APP_NAME")), fmt.Sprintf("%s", os.Getenv("CARTAPP_TOPORT")), 3)
}

func (s *Suite) TearDownSuite() {
	s.ctx.Done()
}

func (s *Suite) TestDeleteEmptyCart() {
	// Тестируем удаление несуществующей корзины.
	// Попытка удалить корзину которая не сущестует.
	httpCode, err := s.cart.DelCart(1)
	require.NoError(s.T(), err, "Ошибка при удалении корзины")
	require.Equal(s.T(), http.StatusNoContent, httpCode, "Код ответа от сервиса не совпадает (удаление корзины)")
}

func (s *Suite) TestUnknownSkuID() {
	// Тестируем удаление несуществующей корзины.
	// Попытка добавить в корзину несуществующий товар.
	itemError := &model.AddItem{}
	itemError.SkuID = 1
	itemError.Count = 1
	itemError.UserID = 2

	httpCode, err := s.cart.AddItem(itemError)
	require.NoError(s.T(), err, "Неверная обработка при попытке добавить некорректный товар в корзину")
	require.Equal(s.T(), http.StatusPreconditionFailed, httpCode, "Код ответа от сервиса не совпадает (добавление несуществующего товара)")
}

func (s *Suite) TestNormalSkuID() {
	// Добавление в корзину корректного товара.
	itemOk := &model.AddItem{}
	itemOk.SkuID = 2958025
	itemOk.Count = 1
	itemOk.UserID = 3

	httpCode, err := s.cart.AddItem(itemOk)
	require.NoError(s.T(), err, "Неверная обработка при попытке добавить верный товар в корзину")
	require.Equal(s.T(), http.StatusOK, httpCode, "Код ответа от сервиса не совпадает (добавление существующего товара)")
}

func (s *Suite) TestNormalSkuIDLowCountInStock() {
	// Добавление в корзину корректного товара, но с большим числом (отсутствует в остатках).
	itemOk := &model.AddItem{}
	itemOk.SkuID = 2958025
	itemOk.Count = 10000
	itemOk.UserID = 6

	httpCode, _ := s.cart.AddItem(itemOk)
	require.Equal(s.T(), http.StatusPreconditionFailed, httpCode, "Код ответа от сервиса не совпадает (добавление существующего товара)")
}

func (s *Suite) TestDelOnePosition() {
	// Тестируем корректность удаления 1 позиции из корзины.
	// Добавление в корзину корректного товара.
	itemOk := &model.AddItem{}
	itemOk.SkuID = 2958025
	itemOk.Count = 1
	itemOk.UserID = 4

	httpCode, err := s.cart.AddItem(itemOk)
	require.NoError(s.T(), err, "Неверная обработка при попытке добавить верный товар в корзину")
	require.Equal(s.T(), http.StatusOK, httpCode, "Код ответа от сервиса не совпадает (добавление существующего товара)")

	// Добавление в корзину второго корректного товара.
	itemOk.SkuID = 773297411

	httpCode, err = s.cart.AddItem(itemOk)
	require.NoError(s.T(), err, "Неверная обработка при попытке добавить верный товар в корзину")
	require.Equal(s.T(), http.StatusOK, httpCode, "Код ответа от сервиса не совпадает (добавление существующего товара)")

	// Получение корзины пользователя (проверяем наличие двух позиций по разным SkuID).
	cart, httpCode, err := s.cart.GetAllUserItems(4)
	require.NoError(s.T(), err, "Неверная обработка при попытке получить корзину пользователя")
	require.Equal(s.T(), http.StatusOK, httpCode, "Код ответа от сервиса не совпадает (получение корзины)")
	require.Equal(s.T(), 2, len(cart.Items), "Число позиций в корзине не совпадает (получение корзины)")

	// Удаление позиции товара из корзины.
	itemDelete := &model.DelItem{}
	itemDelete.SkuID = 773297411
	itemDelete.UserID = 4

	httpCode, err = s.cart.DelItem(itemDelete)
	require.NoError(s.T(), err, "Ошибка при удалении товара из корзины")
	require.Equal(s.T(), http.StatusNoContent, httpCode, "Код ответа от сервиса не совпадает (удаление итема)")

	// Получение корзины пользователя (проверяем наличие одной позиции).
	cart, httpCode, err = s.cart.GetAllUserItems(4)
	require.NoError(s.T(), err, "Неверная обработка при попытке получить корзину пользователя")
	require.Equal(s.T(), http.StatusOK, httpCode, "Код ответа от сервиса не совпадает (получение корзины)")
	require.Equal(s.T(), 1, len(cart.Items), "Число позиций в корзине не совпадает (получение корзины)")
}

func (s *Suite) TestEmptyAfterDelItem() {
	// Тестируем корректность удаления 1 позиции из корзины.
	// Добавление в корзину корректного товара.
	itemOk := &model.AddItem{}
	itemOk.SkuID = 2958025
	itemOk.Count = 1
	itemOk.UserID = 5

	httpCode, err := s.cart.AddItem(itemOk)
	require.NoError(s.T(), err, "Неверная обработка при попытке добавить верный товар в корзину")
	require.Equal(s.T(), http.StatusOK, httpCode, "Код ответа от сервиса не совпадает (добавление существующего товара)")

	// Получение корзины пользователя (проверяем наличие одной позиции).
	cart, httpCode, err := s.cart.GetAllUserItems(5)
	require.NoError(s.T(), err, "Неверная обработка при попытке получить корзину пользователя")
	require.Equal(s.T(), http.StatusOK, httpCode, "Код ответа от сервиса не совпадает (получение корзины)")
	require.Equal(s.T(), 1, len(cart.Items), "Число позиций в корзине не совпадает (получение корзины)")

	// Удаление позиции товара из корзины.
	itemDelete := &model.DelItem{}
	itemDelete.SkuID = 2958025
	itemDelete.UserID = 5

	httpCode, err = s.cart.DelItem(itemDelete)
	require.NoError(s.T(), err, "Ошибка при удалении товара из корзины")
	require.Equal(s.T(), http.StatusNoContent, httpCode, "Код ответа от сервиса не совпадает (удаление итема)")

	// Получение корзины пользователя (проверяем наличие ошибки, что корзины нет).
	_, httpCode, err = s.cart.GetAllUserItems(5)
	require.NoError(s.T(), err, "Ошибка при получении корзины")
	require.Equal(s.T(), http.StatusNotFound, httpCode, "Код ответа от сервиса о при получении пустой корзины не совпадает (получение пустой корзины)")
}

func (s *Suite) TestCheckoutEmptyCart() {
	// Тестируем попытку создать заказ с пустой корзиной.
	_, httpCode, err := s.cart.Checkout(7)
	require.NoError(s.T(), err, "Ошибка при попытке создать заказ по пустой корзине")
	require.Equal(s.T(), http.StatusBadRequest, httpCode, "Код ответа от сервиса не совпадает (заказ по пустой корзине)")
}

func (s *Suite) TestNormalCkeckout() {
	// Чекаут корректной корзины.
	itemOk := &model.AddItem{}
	itemOk.SkuID = 2958025
	itemOk.Count = 1
	itemOk.UserID = 8

	httpCode, err := s.cart.AddItem(itemOk)
	require.NoError(s.T(), err, "Неверная обработка при попытке добавить верный товар в корзину")
	require.Equal(s.T(), http.StatusOK, httpCode, "Код ответа от сервиса не совпадает (добавление существующего товара)")

	_, httpCode, err = s.cart.Checkout(itemOk.UserID)
	require.NoError(s.T(), err, "Ошибка при попытке создать заказ по корректной корзине")
	require.Equal(s.T(), http.StatusOK, httpCode, "Код ответа от сервиса не совпадает (заказ по корректной корзине)")
}
