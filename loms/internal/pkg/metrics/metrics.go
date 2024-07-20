package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/config"
)

var requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Subsystem: config.AppName,
	Name:      "requests_total",
	Help:      "Общее количество запросов, сделанных к сервису. Пример: rate(" + config.AppName + "_requests_total[1m])",
}, []string{"handler"})

var responseCode = promauto.NewCounterVec(prometheus.CounterOpts{
	Subsystem: config.AppName,
	Name:      "response_code",
	Help:      "Коды ответа сервиса. Пример: rate(" + config.AppName + "_response_code[1m])",
}, []string{"handler", "code"})

var responseTime = promauto.NewHistogram(prometheus.HistogramOpts{
	Subsystem: config.AppName,
	Name:      "response_time",
	Buckets:   prometheus.DefBuckets,
	Help:      "Время ответа от сервиса. Пример: rate(" + config.AppName + "_response_time[1m])",
})

var ordersCreated = promauto.NewCounter(prometheus.CounterOpts{
	Subsystem: config.AppName,
	Name:      "orders_created",
	Help:      "Число созданных заказов. Пример: rate(" + config.AppName + "orders_created[1m])",
})

var ordersCreatedError = promauto.NewCounter(prometheus.CounterOpts{
	Subsystem: config.AppName,
	Name:      "orders_create_error",
	Help:      "Число ошибок при создании заказов. Пример: rate(" + config.AppName + "orders_create_error[1m])",
})

var orderStatusChanged = promauto.NewCounterVec(prometheus.CounterOpts{
	Subsystem: config.AppName,
	Name:      "order_status_changed",
	Help:      "Число заказов по которым была смена статуса. Пример: rate(" + config.AppName + "order_status_changed[1m])",
}, []string{"from", "to"})

var externalRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Subsystem: config.AppName,
	Name:      "external_requests_total",
	Help:      "Общее число обращений во внешние сервисы. Пример: rate(" + config.AppName + "external_requests_total[1m])",
}, []string{"service", "handler"})

var externalResponseCode = promauto.NewCounterVec(prometheus.CounterOpts{
	Subsystem: config.AppName,
	Name:      "external_response_code",
	Help:      "Коды ответов от внешних сервисов. Пример: rate(" + config.AppName + "external_response_code[1m])",
}, []string{"service", "handler", "code"})

var externalResponseTime = promauto.NewHistogram(prometheus.HistogramOpts{
	Subsystem: config.AppName,
	Name:      "external_response_time",
	Buckets:   prometheus.DefBuckets,
	Help:      "Время ответа от внешних сервисов. Пример: rate(" + config.AppName + "external_response_time[1m])",
})

var databaseRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Subsystem: config.AppName,
	Name:      "database_requests_total",
	Help:      "Общее число запросов к БД. Пример: rate(" + config.AppName + "database_requests_total[1m])",
}, []string{"repository", "method", "operation"})

var databaseResponseCode = promauto.NewCounterVec(prometheus.CounterOpts{
	Subsystem: config.AppName,
	Name:      "database_response_code",
	Help:      "Статусы ответов от БД. Пример: rate(" + config.AppName + "external_response_code[1m])",
}, []string{"repository", "method", "operation", "code"})

var databaseResponseTime = promauto.NewHistogram(prometheus.HistogramOpts{
	Subsystem: config.AppName,
	Name:      "database_requests_time",
	Buckets:   prometheus.DefBuckets,
	Help:      "Время ответа от БД по запросам. Пример: rate(" + config.AppName + "database_response_time[1m])",
})

// UpdateRequestsTotal метрика суммарного числа запросов к обработчику сервиса.
func UpdateRequestsTotal(handler string) {
	requestsTotal.WithLabelValues(handler).Inc()
}

// UpdateResponseCode распределение кодов ответов по запросам.
func UpdateResponseCode(handler, code string) {
	responseCode.WithLabelValues(handler, code).Inc()
}

// UpdateResponseTime распределение времени ответов по запросам.
func UpdateResponseTime(start time.Time) {
	responseTime.Observe(time.Since(start).Seconds())
}

// UpdateOrdersCreated число созданных заказов.
func UpdateOrdersCreated() {
	ordersCreated.Inc()
}

// UpdateOrdersCreatedError число ошибок при создании заказа.
func UpdateOrdersCreatedError() {
	ordersCreatedError.Inc()
}

// UpdateOrderStatusChanged число заказов по которым была смена статуса.
func UpdateOrderStatusChanged(from, to string) {
	orderStatusChanged.WithLabelValues(from, to).Inc()
}

// UpdateExternalRequestsTotal число обращений во внешние сервисы.
func UpdateExternalRequestsTotal(service, handler string) {
	externalRequestsTotal.WithLabelValues(service, handler).Inc()
}

// UpdateExternalResponseCode распределение кодов ответов от внешних сервисов.
func UpdateExternalResponseCode(service, handler, code string) {
	externalResponseCode.WithLabelValues(service, handler, code).Inc()
}

// UpdateExternalResponseTime распределение времени ответов от внешних сервисов.
func UpdateExternalResponseTime(start time.Time) {
	externalResponseTime.Observe(time.Since(start).Seconds())
}

// UpdateDatabaseRequestsTotal число запросов к БД.
func UpdateDatabaseRequestsTotal(repository, method, operation string) {
	databaseRequestsTotal.WithLabelValues(repository, method, operation).Inc()
}

// UpdateDatabaseResponseCode распределение кодов ответов от БД.
func UpdateDatabaseResponseCode(repository, method, operation, code string) {
	databaseResponseCode.WithLabelValues(repository, method, operation, code).Inc()
}

// UpdateDatabaseResponseTime распределение времени ответов от БД.
func UpdateDatabaseResponseTime(start time.Time) {
	databaseResponseTime.Observe(time.Since(start).Seconds())
}
