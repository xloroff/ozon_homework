//go:build integration

package integrationtest

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/order_store"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/stock_store"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/pkg/db"
)

const (
	user1 = int64(1)
	user2 = int64(2)
	sku1  = int64(773297411)
	sku2  = int64(4678816)
	sku3  = int64(31957466)
	sku4  = int64(28349359)
)

type Suite struct {
	suite.Suite
	ctx          context.Context
	dbClient     db.ClientBD
	logger       logger.ILog
	stockStorage stockstore.Storage
	orderStorage orderstore.Storage
}

func TestIntegrationTest(t *testing.T) {
	suite.Run(t, &Suite{})

}

func (s *Suite) SetupSuite() {
	s.ctx = context.Background()
	s.logger = logger.InitializeLogger("", 1)

	var err error

	TestDB1Url := os.Getenv("DB_NODE_1_CON")
	TestDB2Url := os.Getenv("DB_SYNC_1_CON")
	migratFldr := os.Getenv("MIGRATION_FOLDER")

	err = db.MigrationPool(s.ctx, s.logger, migratFldr, TestDB1Url)
	if err != nil {
		s.Require().NoError(err, "Ошибка миграции БД")
	}

	s.dbClient, err = db.NewClient(s.ctx, TestDB1Url, TestDB2Url)
	if err != nil {
		s.Require().NoError(err, "Ошибка подключения к БД")
	}

	s.stockStorage, err = stockstore.NewReserveStorage(s.ctx, s.logger, s.dbClient)
	if err != nil {
		s.Require().NoError(err, "Ошибка создания хранилища резервов")
	}

	s.orderStorage, err = orderstore.NewOrderStorage(s.ctx, s.logger, s.dbClient)
	if err != nil {
		s.Require().NoError(err, "Ошибка создания хранилища заказов")
	}

}

func (s *Suite) TearDownSuite() {
	s.dbClient.Close()
}

func (s *Suite) TestUnknownSkuIDStockStorage() {
	// Запрос несуществующего SKU в остатках.
	_, err := s.stockStorage.GetAvailableForReserve(123)
	s.Require().Error(err, "Нет ошибки при запросе несуществующего запаса")
}

func (s *Suite) TestReserveStockStorage() {
	// Проверяем корректность резервирования товаров.
	// Получаем количество до выполнения.
	cnt, err := s.stockStorage.GetAvailableForReserve(sku1)
	s.Require().NoError(err, "Ошибка получения количества")
	s.Require().Equal(uint16(60), cnt, "Не совпало количество при старте теста")

	// Резервируем 5.
	err = s.stockStorage.AddReserve(model.AllNeedReserve{
		&model.NeedReserve{
			Sku:   sku1,
			Count: 5,
		},
	})
	s.Require().NoError(err, "Ошибка при резервировании")

	// Получаем количество после резервирования.
	cntTwo, err := s.stockStorage.GetAvailableForReserve(sku1)
	s.Require().NoError(err, "Ошибка при получении количества")
	s.Require().Equal(uint16(55), cntTwo, "Не совпало количество после первого резервирования")

	// Резервируем 30 штук.
	err = s.stockStorage.AddReserve(model.AllNeedReserve{
		&model.NeedReserve{
			Sku:   sku1,
			Count: 30,
		},
	})
	s.Require().NoError(err, "Ошибка при резервировании")

	// Получаем количество после резервирования.
	cntThree, err := s.stockStorage.GetAvailableForReserve(sku1)
	s.Require().NoError(err, "Ошибка при получении количества")
	s.Require().Equal(uint16(25), cntThree, "Не совпало количество после второго резервирования")
}

func (s *Suite) TestCancelReserveStockStorage() {
	// Проверяем корректность снятия резерва
	// Получаем количество до выполнения.
	cnt, err := s.stockStorage.GetAvailableForReserve(sku2)
	s.Require().NoError(err, "Ошибка получения количества")
	s.Require().Equal(uint16(90), cnt, "Не совпало количество при старте теста")

	// Резервируем 5.
	err = s.stockStorage.AddReserve(model.AllNeedReserve{
		&model.NeedReserve{
			Sku:   sku2,
			Count: 5,
		},
	})
	s.Require().NoError(err, "Ошибка при резервировании")

	// Получаем количество после резервирования.
	cntTwo, err := s.stockStorage.GetAvailableForReserve(sku2)
	s.Require().NoError(err, "Ошибка при получении количества")
	s.Require().Equal(uint16(85), cntTwo, "Не совпало количество после первого резервирования")

	// Отменяем резерв
	err = s.stockStorage.CancelReserve(model.AllNeedReserve{
		&model.NeedReserve{
			Sku:   sku2,
			Count: 15,
		},
	})
	s.Require().NoError(err, "Ошибка при отмене резерва")

	// Получаем количество после отмены резервировани.
	cntThree, err := s.stockStorage.GetAvailableForReserve(sku2)
	s.Require().NoError(err, "Ошибка получения количества")
	s.Require().Equal(uint16(100), cntThree, "Не совпало количество после отмены резерва")
}

func (s *Suite) TestCancelTooMahStockStorage() {
	// Отменяем резерв
	err := s.stockStorage.CancelReserve(model.AllNeedReserve{
		&model.NeedReserve{
			Sku:   sku2,
			Count: 10000,
		},
	})
	s.Require().Error(err, "Не получена ошибка при отмене слишком большого количества")
}

func (s *Suite) TestOrdersCreateOrderStorage() {
	// Тестируем создание заказа и получение заказа.
	orderID, err := s.orderStorage.AddOrder(user1, model.OrderItems{
		&model.OrderItem{
			Sku:   sku3,
			Count: 10,
		},
	})
	s.Require().NoError(err, "Ошибка при создании заказа")
	s.Require().GreaterOrEqual(orderID, int64(1), "Не совпал ID созданного заказа")

	// Получаем созданный заказ
	order, err := s.orderStorage.GetOrder(orderID)
	s.Require().NoError(err, "Ошибка при получении созданного заказа")
	s.Require().Equal(user1, order.User, "Не совпал ID пользователя")
	s.Require().Equal(model.OrderStatusNew, order.Status, "Не совпал статус заказа")

	s.Require().Len(order.Items, 1, "Не совпало количество позиций в заказе")
}

func (s *Suite) TestOrdersChangeStatusOrderStorage() {
	// Тестируем создание заказа и смену статуса после его отмены.
	orderID, err := s.orderStorage.AddOrder(user2, model.OrderItems{
		&model.OrderItem{
			Sku:   sku4,
			Count: 10,
		},
	})
	s.Require().NoError(err, "Ошибка при создании заказа")
	s.Require().GreaterOrEqual(orderID, int64(1), "Не совпал ID созданного заказа")

	// Изменяем статусу заказа
	err = s.orderStorage.SetStatus(orderID, model.OrderStatusAwaitingPayment)
	s.Require().NoError(err, "Ошибка при изменении статуса заказа")

	// Получаем созданный заказ
	order, err := s.orderStorage.GetOrder(orderID)
	s.Require().NoError(err, "Ошибка при получении созданного заказа")
	s.Require().Equal(user2, order.User, "Не совпал ID пользователя")
	s.Require().Equal(model.OrderStatusAwaitingPayment, order.Status, "Не совпал статус заказа")
	s.Require().Len(order.Items, 1, "Не совпало количество позиций в заказе")
}
