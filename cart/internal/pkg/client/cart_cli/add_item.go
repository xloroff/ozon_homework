package cartcli

import (
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model"
)

// AddItem добавляет товар в память.
func (cc *Client) AddItem(item *model.AddItem) (code int, err error) {
	req := model.AddItemBody{
		Count: item.Count,
	}

	urlapi := "/user/" + fmt.Sprintf("%d", item.UserID) + "/cart/" + fmt.Sprintf("%d", item.SkuID)

	resp, err := cc.connector.R().
		SetBody(req).
		EnableTrace().
		ForceContentType("application/json").
		Post(urlapi)
	if err != nil {
		return 0, fmt.Errorf("Ошибка отправки запроса в сервис корзины - %w", err)
	}

	return resp.StatusCode(), nil
}
