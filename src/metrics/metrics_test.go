package metrics

import (
	"fmt"
	"github.com/geometry-labs/go-service-template/config"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsApiStart(t *testing.T) {
	assert := assert.New(t)

	// Set env
	os.Setenv("METRICS_PORT", "8888")
	os.Setenv("METRICS_PREFIX", "/metrics")

	config.GetEnvironment()

	// Start metrics server
	MetricsApiStart()

	Metrics["requests_amount"].Inc()
	Metrics["kafka_messages_consumed"].Inc()
	Metrics["websockets_connected"].Inc()
	Metrics["websockets_bytes_written"].Inc()

	resp, err := http.Get(fmt.Sprintf("http://localhost:%s%s", config.Config.MetricsPort, config.Config.MetricsPrefix))
	assert.Equal(nil, err)
	assert.Equal(200, resp.StatusCode)
}

// ISSUE #42
//func TestMetricsWorkerStart(t *testing.T) {
//	assert := assert.New(t)
//
//	// Set env
//	os.Setenv("METRICS_PORT", "8888")
//	os.Setenv("METRICS_PREFIX", "/metrics")
//
//	GetEnvironment()
//
//	// Start metrics server
//	MetricsWorkerStart()
//
//	Metrics["kafka_messages_consumed"].Inc()
//	Metrics["kafka_messages_produced"].Inc()
//
//	resp, err := http.Get(fmt.Sprintf("http://localhost:%s%s", Config.MetricsPort, Config.MetricsPrefix))
//	assert.Equal(nil, err)
//	assert.Equal(200, resp.StatusCode)
//}
