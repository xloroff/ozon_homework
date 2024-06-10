package cartcli

import (
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/e2e/internal/model"
)

// Checkout получает все товары хранящиеся в корзине пользователя.
func (cc *Client) Checkout(userID int64) (*model.OrderCart, int, error) {
	urlapi := "/user/" + fmt.Sprintf("%d", userID) + "/checkout"

	result := &model.OrderCart{}

	resp, err := cc.connector.R().
		EnableTrace().
		ForceContentType("application/json").
		SetResult(result).
		Post(urlapi)
	if err != nil {
		return nil, 0, fmt.Errorf("Ошибка отправки запроса в сервис корзины - %w", err)
	}

	return result, resp.StatusCode(), nil
}
