package orderapi

import (
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/pkg/api/order/v1"
)

// reqItemstoItems конвертирует данные по пачке резервируемых итемов протобафа в подходящие данные под сервис.
func reqItemstoItems(items []*order.OrderCreateRequest_Item) model.AllNeedReserve {
	res := make(model.AllNeedReserve, 0, len(items))

	for _, item := range items {
		res = append(res, convertItem(item))
	}

	return res
}

// convertItem конвертирует данные по одному итему от протобафа в подходящие данные под сервис.
func convertItem(item *order.OrderCreateRequest_Item) *model.NeedReserve {
	return &model.NeedReserve{
		Sku:   item.GetSku(),
		Count: uint16(item.GetCount()),
	}
}

// orderToResponse конвертирует данные по заказу для ответа в протобоф.
func orderToResponse(orderID int64, o *model.Order) *order.Order {
	return &order.Order{
		Id:     orderID,
		User:   o.User,
		Status: o.Status,
		Items:  orderToRespItems(o.Items),
	}
}

// orderToRespItems конвертирует все итемы заказа для ответа в протобоф.
func orderToRespItems(items model.OrderItems) []*order.OrderItem {
	res := make([]*order.OrderItem, 0, len(items))

	for _, item := range items {
		res = append(res, orderToRespItem(item))
	}

	return res
}

// orderToRespItem конвертирует итем заказа для ответа в протобоф.
func orderToRespItem(item *model.OrderItem) *order.OrderItem {
	return &order.OrderItem{
		Sku:   item.Sku,
		Count: uint64(item.Count),
	}
}
