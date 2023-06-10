package main

import (
	"github.com/geometry-labs/go-service-template/config"
	"github.com/geometry-labs/go-service-template/global"
	"github.com/geometry-labs/go-service-template/kafka"
	"github.com/geometry-labs/go-service-template/logging"
	"github.com/geometry-labs/go-service-template/metrics"
	"github.com/geometry-labs/go-service-template/worker/loader"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"

	"github.com/geometry-labs/go-service-template/worker/transformers"
)

func main() {

	config.GetEnvironment()

	logging.StartLoggingInit()
	zap.S().Debug("Main: Starting logging with level ", config.Config.LogLevel)

	// Start Prometheus client
	metrics.MetricsWorkerStart()

	// Start Health server
	//healthcheck.Start()

	// Start kafka consumer
	kafka.StartWorkerConsumers()

	// Start kafka Producer
	kafka.StartProducers()
	// Wait for Kafka
	//time.Sleep(1 * time.Second)

	// Start Postgres loader
	loader.StartBlockRawsLoader()

	// Start transformers
	transformers.StartBlocksTransformer()

	// Listen for close sig
	// Register for interupt (Ctrl+C) and SIGTERM (docker)

	//create a notification channel to shutdown
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		zap.S().Info("Shutting down...")
		global.ShutdownChan <- 1
	}()

	<-global.ShutdownChan
}
