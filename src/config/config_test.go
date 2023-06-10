package config

import (
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
	assert := assert.New(t)

	// Set env
	env_map := map[string]string{
		"VERSION":                 "version",
		"NAME":                    "name",
		"PORT":                    "port",
		"HEALTH_PORT":             "health_port",
		"METRICS_PORT":            "metrics_port",
		"REST_PREFIX":             "rest_prefix",
		"WEBSOCKET_PREFIX":        "websocket_prefix",
		"HEALTH_PREFIX":           "health_prefix",
		"METRICS_PREFIX":          "metrics_prefix",
		"HEALTH_POLLING_INTERVAL": "5",
		"LOG_LEVEL":               "log_level",
		"LOG_TO_FILE":             "true",
		"NETWORK_NAME":            "network_name",
		"KAFKA_BROKER_URL":        "kafka_broker_url",
		"SCHEMA_REGISTRY_URL":     "schema_registry_url",
		"KAFKA_GROUP_ID":          "kafka_group_id",
		"CONSUMER_TOPICS":         "topic_1,topic_2",
		"PRODUCER_TOPICS":         "topic_1,topic_2,topic_3",
		"SCHEMA_NAMES":            "schema_1:schema_1,schema_2:schema_2",
	}

	for k, v := range env_map {
		os.Setenv(k, v)
	}

	// Load env
	GetEnvironment()

	// Check env
	assert.Equal(env_map["VERSION"], Config.Version)
	assert.Equal(env_map["NAME"], Config.Name)
	assert.Equal(env_map["PORT"], Config.Port)
	assert.Equal(env_map["HEALTH_PORT"], Config.HealthPort)
	assert.Equal(env_map["METRICS_PORT"], Config.MetricsPort)
	assert.Equal(env_map["REST_PREFIX"], Config.RestPrefix)
	assert.Equal(env_map["WEBSOCKET_PREFIX"], Config.WebsocketPrefix)
	assert.Equal(env_map["HEALTH_PREFIX"], Config.HealthPrefix)
	assert.Equal(env_map["METRICS_PREFIX"], Config.MetricsPrefix)
	assert.Equal(5, Config.HealthPollingInterval)
	assert.Equal(env_map["LOG_LEVEL"], Config.LogLevel)
	assert.Equal(true, Config.LogToFile)
	assert.Equal(env_map["NETWORK_NAME"], Config.NetworkName)
	assert.Equal(env_map["KAFKA_BROKER_URL"], Config.KafkaBrokerURL)
	assert.Equal(env_map["SCHEMA_REGISTRY_URL"], Config.SchemaRegistryURL)
	assert.Equal(env_map["KAFKA_GROUP_ID"], Config.KafkaGroupID)
	assert.Equal(2, len(Config.Topics.ConsumerTopics))
	assert.Equal(3, len(Config.Topics.ProducerTopics))
	assert.Equal("schema_1", Config.Topics.SchemaNames["schema_1"])
	assert.Equal("schema_2", Config.Topics.SchemaNames["schema_2"])
}
