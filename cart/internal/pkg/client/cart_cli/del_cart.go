package cartcli

import (
	"fmt"
)

// DelCart удаляет корзину пользователя в сервисе корзины.
func (cc *Client) DelCart(userID int64) (code int, err error) {
	urlapi := "/user/" + fmt.Sprintf("%d", userID) + "/cart"

	resp, err := cc.connector.R().
		EnableTrace().
		ForceContentType("application/json").
		Delete(urlapi)
	if err != nil {
		return 0, fmt.Errorf("Ошибка отправки запроса в сервис корзины - %w", err)
	}

	return resp.StatusCode(), nil
}
