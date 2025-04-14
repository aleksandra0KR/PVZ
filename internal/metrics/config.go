package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "endpoint", "status"},
	)

	ResponseDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_duration_seconds",
			Help:    "Histogram of response durations for HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	PvzCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pvz_created_total",
			Help: "Total number of PVZs created.",
		},
	)

	ReceptionsCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "receptions_created_total",
			Help: "Total number of receptions created.",
		},
	)

	ProductsAdded = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "products_added_total",
			Help: "Total number of products added.",
		},
	)
)

func Init() {
	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(ResponseDuration)
	prometheus.MustRegister(PvzCreated)
	prometheus.MustRegister(ReceptionsCreated)
	prometheus.MustRegister(ProductsAdded)
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()

		status := c.Writer.Status()
		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = "unknown"
		}

		RequestsTotal.WithLabelValues(c.Request.Method, endpoint, http.StatusText(status)).Inc()
		ResponseDuration.WithLabelValues(c.Request.Method, endpoint).Observe(duration)
	}
}

func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
