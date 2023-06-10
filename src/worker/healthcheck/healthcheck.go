package healthcheck

import (
	"github.com/geometry-labs/go-service-template/config"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"time"

	"github.com/InVisionApp/go-health/v2"
	"github.com/InVisionApp/go-health/v2/checkers"
	"github.com/InVisionApp/go-health/v2/handlers"
)

// TODO split API and WORKER
func Start() {
	// Create a new health instance
	h := health.New()

	// Create a couple of checks
	blocksCheckerURL, _ := url.Parse("http://localhost:" + config.Config.Port + "/version")
	blocksChecker, _ := checkers.NewHTTP(&checkers.HTTPConfig{
		URL: blocksCheckerURL,
	})

	// Add the checks to the health instance
	h.AddChecks([]*health.Config{
		{
			Name:     "blocks-rest-check",
			Checker:  blocksChecker,
			Interval: time.Duration(config.Config.HealthPollingInterval) * time.Second,
			Fatal:    true,
		},
	})

	//  Start the healthcheck process
	if err := h.Start(); err != nil {
		zap.S().Fatalf("Unable to start healthcheck: %v", err)
	}

	// Define a healthcheck endpoint and use the built-in JSON handler
	http.HandleFunc(config.Config.HealthPrefix, handlers.NewJSONHandlerFunc(h, nil))
	go http.ListenAndServe(":"+config.Config.HealthPort, nil)
	zap.S().Info("Started Healthcheck:", config.Config.HealthPort)
}
