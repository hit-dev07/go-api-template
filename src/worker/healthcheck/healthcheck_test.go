package healthcheck

import (
	"github.com/geometry-labs/go-service-template/config"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/geometry-labs/go-service-template/api/routes"
)

func init() {
	//core.GetEnvironment()
	config.Vars.ConfigFile = "config.api.test"
	config.Vars.ConfigType = "yaml"
	config.Vars.ConfigPath = "../../../envfiles"
	config.ConfigInit()
}

func TestHealthCheck(t *testing.T) {
	assert := assert.New(t)

	// Start api
	routes.Start()

	// Start healthcheck
	Start()

	resp, err := http.Get("http://localhost:" + config.Config.HealthPort + config.Config.HealthPrefix)
	assert.Equal(nil, err)
	assert.Equal(200, resp.StatusCode)
}
