package middleware

import (
	"net/http"
	"time"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/metrics"
)

// Metrics общие метрики доступа к прикладам.
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.UpdateRequestsTotal()
		defer metrics.UpdateResponseTime(time.Now().UTC())

		next.ServeHTTP(w, r)
	})
}
