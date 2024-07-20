package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/config"
)

var requestsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Subsystem: config.AppName,
	Name:      "requests_total",
	Help:      "Общее количество запросов, сделанных к сервису. Пример: rate(" + config.AppName + "_requests_total[1m])",
})

var requestsTotalURL = promauto.NewCounterVec(prometheus.CounterOpts{
	Subsystem: config.AppName,
	Name:      "requests_total_url",
	Help:      "Общее количество запросов, сделанных к сервису. Пример: rate(" + config.AppName + "requests_total_url[1m])",
}, []string{"method", "url"})

var responseCode = promauto.NewCounterVec(prometheus.CounterOpts{
	Subsystem: config.AppName,
	Name:      "response_code",
	Help:      "Коды ответа сервиса. Пример: rate(" + config.AppName + "response_code[1m])",
}, []string{"method", "url", "code"})

var responseTime = promauto.NewHistogram(prometheus.HistogramOpts{
	Subsystem: config.AppName,
	Name:      "response_time",
	Buckets:   prometheus.DefBuckets,
	Help:      "Время ответа от сервиса. Пример: rate(" + config.AppName + "response_time[1m])",
})

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
	Help:      "Общее число запросов к хранилищу памяти. Пример: rate(" + config.AppName + "database_requests_total[1m])",
}, []string{"repository", "method"})

var databaseResponseCode = promauto.NewCounterVec(prometheus.CounterOpts{
	Subsystem: config.AppName,
	Name:      "database_response_code",
	Help:      "Статусы ответов от хранилища памяти. Пример: rate(" + config.AppName + "external_response_code[1m])",
}, []string{"repository", "method", "operation", "code"})

var databaseResponseTime = promauto.NewHistogram(prometheus.HistogramOpts{
	Subsystem: config.AppName,
	Name:      "database_requests_time",
	Buckets:   prometheus.DefBuckets,
	Help:      "Время ответа от хранилища памяти по запросам. Пример: rate(" + config.AppName + "database_response_time[1m])",
})

var inMemoryItemCount = promauto.NewGauge(prometheus.GaugeOpts{
	Subsystem: config.AppName,
	Name:      "in_memory_item_count",
	Help:      "Текущее количество элементов в хранилище памяти.",
})

// UpdateRequestsTotal метрика суммарного числа запросов к обработчику сервиса.
func UpdateRequestsTotal() {
	requestsTotal.Inc()
}

// UpdateRequestsTotalWithURL метрика суммарного числа запросов к обработчику сервиса в разделении по URL.
func UpdateRequestsTotalWithURL(method, url string) {
	requestsTotalURL.WithLabelValues(method, url).Inc()
}

// UpdateResponseCode распределение кодов ответов по запросам.
func UpdateResponseCode(method, url, code string) {
	responseCode.WithLabelValues(method, url, code).Inc()
}

// UpdateResponseTime распределение времени ответов по запросам.
func UpdateResponseTime(start time.Time) {
	responseTime.Observe(time.Since(start).Seconds())
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

// UpdateExternalResponseDuration распределение времени ответов от внешних сервисов.
func UpdateExternalResponseDuration(duration time.Duration) {
	externalResponseTime.Observe(duration.Seconds())
}

// UpdateDatabaseRequestsTotal число запросов к хранилищу памяти.
func UpdateDatabaseRequestsTotal(repository, method string) {
	databaseRequestsTotal.WithLabelValues(repository, method).Inc()
}

// UpdateDatabaseResponseCode распределение кодов ответов от хранилища памяти.
func UpdateDatabaseResponseCode(repository, method, operation, code string) {
	databaseResponseCode.WithLabelValues(repository, method, operation, code).Inc()
}

// UpdateDatabaseResponseTime распределение времени ответов от  хранилища памяти.
func UpdateDatabaseResponseTime(start time.Time) {
	databaseResponseTime.Observe(time.Since(start).Seconds())
}

// UpdateInMemoryItemCount число элементов в InMemory хранилище.
func UpdateInMemoryItemCount(count int) {
	inMemoryItemCount.Set(float64(count))
}
