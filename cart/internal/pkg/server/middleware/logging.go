package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

// MiddlewareLogging включаем логирование по всем входящим запросам.
func MiddlewareLogging(ctx context.Context) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Сначала общее инфо.
			executeInfo := []zap.Field{
				zap.String("url_path", r.URL.Path),
			}

			ctx = logger.Set(ctx, executeInfo)
			logger.Infof(ctx, "service_access")

			executeDebug := []zap.Field{
				zap.String("ip", ReadUserIP(r)),
				zap.String("content_type", r.Header.Get("Content-Type")),
				// TODO добавить еще что-то полезное, типа таймингов и пр. но как понимаю рано, будет на сл. ДЗ.
			}

			executeDebug = append(executeInfo, executeDebug...)
			ctx = logger.Set(ctx, executeDebug)
			logger.Debugf(ctx, "service_access")

			h.ServeHTTP(w, r)
		})
	}
}

// ReadUserIP вычисляем по IP.
func ReadUserIP(r *http.Request) string {
	ipAddress := r.Header.Get("X-Real-Ip")

	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}

	return ipAddress
}