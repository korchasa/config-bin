package prometheus

import (
	"strconv"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
)

type Prometheus struct {
	requestsTotal             *prom.CounterVec
	successfulRequestsTotal   *prom.CounterVec
	failedRequestsTotal       *prom.CounterVec
	successfulRequestDuration *prom.HistogramVec
	eventSendsTotal           *prom.CounterVec
	kafkaRequestsTotal        *prom.CounterVec
	kafkaErrorsTotal          *prom.CounterVec
	kafkaRequestsDuration     *prom.HistogramVec
}

func New() *Prometheus {
	buckets := []float64{
		0.001, 0.002, 0.005,
		0.01, 0.02, 0.05,
		0.1, 0.2, 0.5,
		1, 2, 5,
		10, 20, 60,
	}
	metrics := buildProm(buckets)

	prom.MustRegister(metrics.requestsTotal)
	prom.MustRegister(metrics.successfulRequestsTotal)
	prom.MustRegister(metrics.failedRequestsTotal)
	prom.MustRegister(metrics.successfulRequestDuration)
	prom.MustRegister(metrics.eventSendsTotal)
	prom.MustRegister(metrics.kafkaRequestsTotal)
	prom.MustRegister(metrics.kafkaErrorsTotal)
	prom.MustRegister(metrics.kafkaRequestsDuration)

	return metrics
}

func buildProm(buckets []float64) *Prometheus {
	return &Prometheus{
		requestsTotal: prom.NewCounterVec(
			prom.CounterOpts{
				Name: "http_producer_requests_total",
				Help: "Number of requests",
			},
			[]string{},
		),
		successfulRequestsTotal: prom.NewCounterVec(
			prom.CounterOpts{
				Name: "http_producer_successful_requests_total",
				Help: "Number of successful requests",
			},
			[]string{"code"},
		),
		failedRequestsTotal: prom.NewCounterVec(
			prom.CounterOpts{
				Name: "http_producer_failed_requests_total",
				Help: "Number of failed requests.",
			},
			[]string{"code"},
		),
		successfulRequestDuration: prom.NewHistogramVec(
			prom.HistogramOpts{
				Name:    "http_producer_successful_requests_duration",
				Help:    "Duration of successful requests.",
				Buckets: buckets,
			},
			[]string{},
		),
		eventSendsTotal: prom.NewCounterVec(
			prom.CounterOpts{
				Name: "http_producer_event_sends_total",
				Help: "Number of event sends",
			},
			[]string{"type", "req_app_id", "req_app_instance", "req_app_env", "req_app_project", "ack"},
		),
		kafkaRequestsTotal: prom.NewCounterVec(
			prom.CounterOpts{
				Name: "http_producer_kafka_requests_total",
				Help: "Number of kafka requests.",
			},
			[]string{},
		),
		kafkaErrorsTotal: prom.NewCounterVec(
			prom.CounterOpts{
				Name: "http_producer_kafka_errors_total",
				Help: "Number of kafka requests.",
			},
			[]string{},
		),
		kafkaRequestsDuration: prom.NewHistogramVec(
			prom.HistogramOpts{
				Name:    "http_producer_kafka_requests_duration",
				Help:    "Duration of successful requests.",
				Buckets: buckets,
			},
			[]string{},
		),
	}
}

func (p *Prometheus) IncRequestsCount() {
	p.requestsTotal.WithLabelValues().Inc()
}

func (p *Prometheus) IncSuccessfulRequestsCount(httpCode int) {
	p.successfulRequestsTotal.
		WithLabelValues(strconv.Itoa(httpCode)).
		Inc()
}

func (p *Prometheus) IncFailedRequestsCount(httpCode int) {
	p.failedRequestsTotal.
		WithLabelValues(strconv.Itoa(httpCode)).
		Inc()
}

func (p *Prometheus) ObserveSuccessfulRequestDuration(dur time.Duration) {
	p.successfulRequestDuration.WithLabelValues().Observe(dur.Seconds())
}

func (p *Prometheus) IncEventSendsCount(
	eventType string,
	appID string,
	appInstance string,
	appEnv string,
	project string,
	ack bool,
) {
	p.eventSendsTotal.WithLabelValues(eventType, appID, appInstance, appEnv, project, strconv.FormatBool(ack)).Inc()
}

func (p *Prometheus) IncKafkaRequestsCount() {
	p.kafkaRequestsTotal.WithLabelValues().Inc()
}

func (p *Prometheus) IncKafkaErrorsCount() {
	p.kafkaErrorsTotal.WithLabelValues().Inc()
}

func (p *Prometheus) ObserveKafkaRequestDuration(dur time.Duration) {
	p.kafkaRequestsDuration.WithLabelValues().Observe(dur.Seconds())
}
