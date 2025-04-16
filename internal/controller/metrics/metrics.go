package metrics

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
	"timeline/internal/config"
	"timeline/internal/controller/scope"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Prometheus struct {
	server     *http.Server
	serverOnce sync.Once
	wg         sync.WaitGroup
	log        *zap.Logger
}

func NewPrometheusExporter(cfg config.Prometheus, logger *zap.Logger) Prometheus {
	metricsmux := http.NewServeMux()
	metricsmux.Handle(scope.PathMetrics, promhttp.Handler())
	return Prometheus{
		server: &http.Server{
			Addr:    cfg.Host + ":" + cfg.Port,
			Handler: metricsmux,
		},
		log: logger,
	}
}

func (p *Prometheus) Launch(errch chan error) {
	p.serverOnce.Do(func() {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errch <- fmt.Errorf("prometheus server error: %w", err)
			}
		}()
	})
}

func (p *Prometheus) Shutdown(cancelCtx context.Context, timeout time.Duration) {
	timeoutCtx, cancel := context.WithTimeout(cancelCtx, timeout)
	defer cancel()
	if err := p.server.Shutdown(timeoutCtx); err != nil {
		p.log.Error("failed to shutdown HTTP server", zap.Error(err))
	}
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		p.log.Info("HTTP server is closed successfully")
	case <-timeoutCtx.Done():
		p.log.Error("timeout while closing HTTP server", zap.Error(timeoutCtx.Err()))
	}
}

type Metrics struct {
	RequestDuration *prometheus.HistogramVec
	RequestCounter  *prometheus.CounterVec
}

func NewDefaultMetrics() *Metrics {
	return &Metrics{
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Время обработки HTTP-запроса в секундах",
				Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5}, // Корзины для гистограммы
			},
			[]string{"method", "path", "status"}, // Метки: метод, путь, статус-код
		),
		RequestCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Общее количество HTTP-запросов",
			},
			[]string{"method", "path", "status"}, // Метки: метод, путь, статус-код
		),
	}
}

func (m *Metrics) UpdateRequestMetrics(method, urlpath, status string, duration time.Duration) {
	m.RequestDuration.WithLabelValues(
		method,
		urlpath,
		status,
	).Observe(duration.Seconds())
	m.RequestCounter.WithLabelValues(
		method,
		urlpath,
		status,
	).Inc()
}
