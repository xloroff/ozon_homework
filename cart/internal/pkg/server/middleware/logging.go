package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

// Logging включаем логирование по всем входящим запросам.
func Logging(ctx context.Context, l logger.ILog) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executeDebug := []zap.Field{
				zap.String("ip", ReadUserIP(r)),
				zap.String("content_type", r.Header.Get("Content-Type")),
				// TODO добавить еще что-то полезное, типа таймингов и пр. но как понимаю рано, будет на сл. ДЗ.
			}

			ctx = logger.Set(ctx, executeDebug)
			l.Debugf(ctx, "service_access")

			h.ServeHTTP(w, r)
		})
	}
}

// ReadUserIP вычисляем по IP.
func ReadUserIP(r *http.Request) string {
	ipAddress := r.Header.Get("X-Real-Ip")
	if ipAddress != "" {
		return ipAddress
	}

	ipAddress = r.Header.Get("X-Forwarded-For")
	if ipAddress != "" {
		return ipAddress
	}

	return r.RemoteAddr
}
