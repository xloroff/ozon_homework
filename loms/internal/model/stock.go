package model

// ReserveItem итем и его зарезервированное количество.
type ReserveItem struct {
	Sku        int64  `json:"sku"`
	TotalCount uint16 `json:"total_count"`
	Reserved   uint16 `json:"reserved"`
}

// ReserveItems список всех зарезервированных итемов с количеством.
type ReserveItems []*ReserveItem

// AllReserveItems все резервы заказов(итемов) в памяти.
type AllReserveItems map[int64]*ReserveItem

// NeedReserve запросы на резерв по итемам.
type NeedReserve struct {
	Sku   int64
	Count uint16
}

// AllNeedReserve все резервируемые итемы.
type AllNeedReserve []*NeedReserve
