package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type ConfigStruct struct {
	Name string `mapstructure:"NAME"`

	// Ports
	Port        string `mapstructure:"PORT"`
	HealthPort  string `mapstructure:"HEALTH_PORT"`
	MetricsPort string `mapstructure:"METRICS_PORT"`

	// Prefix
	RestPrefix      string `mapstructure:"REST_PREFIX"`
	WebsocketPrefix string `mapstructure:"WEBSOCKET_PREFIX"`
	HealthPrefix    string `mapstructure:"HEALTH_PREFIX"`
	MetricsPrefix   string `mapstructure:"METRICS_PREFIX"`

	// Monitoring
	HealthPollingInterval int    `mapstructure:"HEALTH_POLLING_INTERVAL"`
	LogLevel              string `mapstructure:"LOG_LEVEL"`
	LogToFile             bool   `mapstructure:"LOG_TO_FILE"`
	NetworkName           string `mapstructure:"NETWORK_NAME"`

	// Kafka
	KafkaBrokerURL    string `mapstructure:"KAFKA_BROKER_URL"`
	SchemaRegistryURL string `mapstructure:"SCHEMA_REGISTRY_URL"`
	KafkaGroupID      string `mapstructure:"KAFKA_GROUP_ID"`

	// Topics
	ConsumerTopics   []string          `mapstructure:"CONSUMER_TOPICS"`
	ProducerTopics   []string          `mapstructure:"PRODUCER_TOPICS"`
	SchemaNameTopics map[string]string `mapstructure:"SCHEMA_NAME_TOPICS"`

	// DB
	DbDriver   string `mapstructure:"DB_DRIVER"`
	DbHost     string `mapstructure:"DB_HOST"`
	DbPort     string `mapstructure:"DB_PORT"`
	DbUser     string `mapstructure:"DB_USER"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbName     string `mapstructure:"DB_DBNAME"`
	DbSslmode  string `mapstructure:"DB_SSLMODE"`
	DbTimezone string `mapstructure:"DB_TIMEZONE"`

	// Mongo
	MongoHost string `mapstructure:"MONGO_HOST"`
	MongoPort string `mapstructure:"MONGO_PORT"`
}

var Config ConfigStruct

// GetEnvironment Run once on main.go
func GetEnvironment() {
	////Get environment variable file
	//env_file := os.Getenv("ENV_FILE")
	//if env_file != "" {
	//	_ = godotenv.Load(env_file)
	//} else {
	//	_ = godotenv.Load()
	//}
	//
	//err := envconfig.Process("", &Vars)
	//if err != nil {
	//	log.Fatalf("ERROR: envconfig - %s\n", err.Error())
	//}

	// Fill Config struct
	config := ConfigInit()
	log.Println("config:", config)
}

func ConfigInit() ConfigStruct {
	//filename := os.Getenv("ENV_FILE")
	filename := "env.worker.test"
	viper.SetConfigName(filename)
	viper.SetConfigType("env")
	viper.AddConfigPath("./envfiles")

	// Set Defaults
	setDefaults()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s, Error: %s \n", filename, err))
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v\n", err))
	}

	return Config
}

func setDefaults() {
	viper.SetDefault("Name", "blocks service")

	viper.SetDefault("Port", "8000")
	viper.SetDefault("HealthPort", "8080")
	viper.SetDefault("MetricsPort", "9400")

	viper.SetDefault("RestPrefix", "/rest")
	viper.SetDefault("WebsocketPrefix", "/ws")
	viper.SetDefault("HealthPrefix", "/healthcheck")
	viper.SetDefault("MetricsPrefix", "/metrics")

	viper.SetDefault("HealthPollingInterval", 10)
	viper.SetDefault("LogLevel", "INFO")
	viper.SetDefault("LogToFile", false)
	viper.SetDefault("NetworkName", "mainnet")

	viper.SetDefault("KafkaBrokerURL", "")
	viper.SetDefault("SchemaRegistryURL", "")
	viper.SetDefault("KafkaGroupID", "websocket-group")

	viper.SetDefault("ConsumerTopics", "[blocks]")
	viper.SetDefault("ProducerTopics", "[blocks-ws]")
	viper.SetDefault("SchemaNameTopics", "blocks:block_raw, blocks-ws:block_raw")
	viper.SetDefault("DbDriver", "postgres")
	viper.SetDefault("DbHost", "localhost")
	viper.SetDefault("DbPort", "5432")
	viper.SetDefault("DbUser", "postgres")
	viper.SetDefault("DbPassword", "changeme")
	viper.SetDefault("DbName", "test_db")
	viper.SetDefault("DbSslmode", "disable")
	viper.SetDefault("DbTimezone", "UTC")

	viper.SetDefault("MongoHost", "localhost")
	viper.SetDefault("MongoPort", "27017")
}
