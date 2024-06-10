package cartcli

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/e2e/internal/model"
)

// CartClient интерфейс клиента для сервиса продуктов.
type CartClient interface {
	AddItem(item *model.AddItem) (int, error)
	GetAllUserItems(userID int64) (*model.FullUserCart, int, error)
	DelItem(item *model.DelItem) (int, error)
	DelCart(userID int64) (int, error)
	Checkout(userID int64) (*model.OrderCart, int, error)
}

// Client структура клиента для сервиса cart.
type Client struct {
	connector *resty.Client
}

// NewCartClient создает клиента для подключения к сервису cart.
func NewCartClient(host, port string, retry int) CartClient {
	connector := resty.New()

	connector.
		// SetRetryCount устанавливает число ретраев.
		SetRetryCount(retry).
		// SetRetryWaitTime таймаут между попытками.
		SetRetryWaitTime(200 * time.Millisecond).
		// SetRetryMaxWaitTime сколько ждуем выполнения попытки.
		SetRetryMaxWaitTime(10 * time.Second).
		// SetTimeout устанавливаем общий таймаут.
		SetTimeout(15 * time.Second).
		// SetRetryAfter что делать если попытки закончились.
		SetRetryAfter(func(_ *resty.Client, _ *resty.Response) (time.Duration, error) {
			return 0, fmt.Errorf("Попытка связи с сервисом продуктов не удалась - %w", connector.Error)
		})

	u, err := url.Parse(host + ":" + port)
	if err != nil {
		panic(err)
	}

	connector.SetBaseURL(u.String())

	return &Client{
		connector: connector,
	}
}
