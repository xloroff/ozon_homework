package cartcli

import (
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model"
)

// GetAllUserItems получает все товары хранящиеся в корзине пользователя.
func (cc *Client) GetAllUserItems(userID int64) (*model.FullUserCart, int, error) {
	urlapi := "/user/" + fmt.Sprintf("%d", userID) + "/cart/list"

	result := &model.FullUserCart{}

	resp, err := cc.connector.R().
		EnableTrace().
		ForceContentType("application/json").
		SetResult(result).
		Get(urlapi)
	if err != nil {
		return nil, 0, fmt.Errorf("Ошибка отправки запроса в сервис корзины - %w", err)
	}

	return result, resp.StatusCode(), nil
}
