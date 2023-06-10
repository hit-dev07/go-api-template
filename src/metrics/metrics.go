package metrics

import (
	"github.com/geometry-labs/go-service-template/config"
	"go.uber.org/zap"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Metrics map[string]prometheus.Counter

func MetricsApiStart() {
	Metrics = make(map[string]prometheus.Counter)

	createApiGauges()

	// Start server
	http.Handle(config.Config.MetricsPrefix, promhttp.Handler())
	go http.ListenAndServe(":"+config.Config.MetricsPort, nil)
	zap.S().Info("Started Metrics:", config.Config.MetricsPort)
}

func MetricsWorkerStart() {
	Metrics = make(map[string]prometheus.Counter)

	// Create gauges
	createWorkerGauges()

	// Start server
	http.Handle(config.Config.MetricsPrefix, promhttp.Handler())
	go http.ListenAndServe(":"+config.Config.MetricsPort, nil)
	zap.S().Info("Started Metrics:", config.Config.MetricsPort)
}

func createApiGauges() {
	Metrics["requests_amount"] = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "requests_amount",
		Help:        "amount of requests",
		ConstLabels: prometheus.Labels{"network_name": config.Config.NetworkName},
	})
	Metrics["kafka_messages_consumed"] = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "kafka_messages_consumed",
		Help:        "amount of messageds from kafka consumed",
		ConstLabels: prometheus.Labels{"network_name": config.Config.NetworkName},
	})
	Metrics["websockets_connected"] = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "websockets_connected",
		Help:        "amount of websockets that have connected to the server",
		ConstLabels: prometheus.Labels{"network_name": config.Config.NetworkName},
	})
	Metrics["websockets_bytes_written"] = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "websockets_bytes_written",
		Help:        "amount of bytes written through websockets",
		ConstLabels: prometheus.Labels{"network_name": config.Config.NetworkName},
	})
}

func createWorkerGauges() {
	Metrics["kafka_messages_consumed"] = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "kafka_messages_consumed",
		Help:        "amount of messages from kafka consumed",
		ConstLabels: prometheus.Labels{"network_name": config.Config.NetworkName},
	})
	Metrics["kafka_messages_produced"] = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "kafka_messages_produced",
		Help:        "amount of messages from kafka produced",
		ConstLabels: prometheus.Labels{"network_name": config.Config.NetworkName},
	})
}
