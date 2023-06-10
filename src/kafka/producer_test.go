package kafka

import (
	"github.com/geometry-labs/go-service-template/config"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/Shopify/sarama.v1"
)

func init() {
	// os.Setenv("LOG_LEVEL", "DEBUG")

	//core.GetEnvironment()
	config.Vars.ConfigFile = "config.api.test"
	config.Vars.ConfigType = "yaml"
	config.Vars.ConfigPath = "../../envfiles"
	config.ConfigInit()
}

func TestKafkaTopicProducer(t *testing.T) {

	topic_name := "mock-topic"

	// Mock broker
	mock_broker_id := int32(1)
	mock_broker := sarama.NewMockBroker(t, mock_broker_id)
	defer mock_broker.Close()

	mock_broker.SetHandlerByMap(map[string]sarama.MockResponse{
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetBroker(mock_broker.Addr(), mock_broker.BrokerID()).
			SetLeader(topic_name, 0, mock_broker.BrokerID()),
		"ProduceRequest": sarama.NewMockProduceResponse(t),
	})

	KafkaTopicProducers[topic_name] = &KafkaTopicProducer{
		mock_broker.Addr(),
		topic_name,
		make(chan *sarama.ProducerMessage),
	}

	go KafkaTopicProducers[topic_name].produceTopic()

	msg_key := "KEY"
	msg_value := "VALUE"
	go func() {
		for {
			KafkaTopicProducers[topic_name].TopicChan <- &sarama.ProducerMessage{
				Topic:     topic_name,
				Partition: -1,
				Key:       sarama.StringEncoder(msg_key),
				Value:     sarama.StringEncoder(msg_value),
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	total_bytes_read := 0
	total_bytes_written := 0
	mock_broker.SetNotifier(func(bytes_read int, bytes_written int) {
		log.Debug("MOCK NOTIFIER: bytes_read=", bytes_read, " bytes_written=", bytes_written)

		total_bytes_read += bytes_read
		total_bytes_written += bytes_written
	})

	loops := 0
	for {
		if loops > 10 {
			t.Fatal("No messages to mock broker")
		}
		if total_bytes_read > 400 && total_bytes_written > 200 {
			// Test passed
			return
		}

		time.Sleep(time.Second)
		loops++
	}
}
