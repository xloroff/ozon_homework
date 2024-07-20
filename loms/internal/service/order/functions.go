package orderservice

import (
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
)

func resevToOrders(items model.AllNeedReserve) model.OrderItems {
	if items == nil {
		return nil
	}

	res := make(model.OrderItems, 0, len(items))
	for _, item := range items {
		res = append(res, reserveToItem(item))
	}

	return res
}

func reserveToItem(item *model.NeedReserve) *model.OrderItem {
	return &model.OrderItem{
		Sku:   item.Sku,
		Count: item.Count,
	}
}

func orderToReserve(items model.OrderItems) model.AllNeedReserve {
	if items == nil {
		return nil
	}

	res := make([]*model.NeedReserve, 0, len(items))
	for _, item := range items {
		res = append(res, orderToItem(item))
	}

	return res
}

func orderToItem(item *model.OrderItem) *model.NeedReserve {
	return &model.NeedReserve{
		Sku:   item.Sku,
		Count: item.Count,
	}
}
