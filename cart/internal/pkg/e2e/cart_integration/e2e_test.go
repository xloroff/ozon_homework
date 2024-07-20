//go:build e2e

package cartintegration

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/app"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/client/cart_cli"
	_ "gitlab.ozon.dev/xloroff/ozon-hw-go/testinginit"
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
	s.cart = cartcli.NewCartClient("http://localhost", fmt.Sprintf("%d", 8082), 3)

	go func(ctx context.Context) {
		if err := app.NewApp(s.ctx).Run(); err != nil {
			require.NoError(s.T(), err, "Приложение не должно падать")
		}

		<-ctx.Done()
	}(s.ctx)
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
