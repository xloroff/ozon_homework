package cartcli

import (
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/e2e/internal/model"
)

// DelItem удаляет товар из корзины пользователя в сервисе корзин.
func (cc *Client) DelItem(item *model.DelItem) (code int, err error) {
	urlapi := "/user/" + fmt.Sprintf("%d", item.UserID) + "/cart/" + fmt.Sprintf("%d", item.SkuID)

	resp, err := cc.connector.R().
		EnableTrace().
		ForceContentType("application/json").
		Delete(urlapi)
	if err != nil {
		return 0, fmt.Errorf("Ошибка отправки запроса в сервис корзины - %w", err)
	}

	return resp.StatusCode(), nil
}
