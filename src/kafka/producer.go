package kafka

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/geometry-labs/go-service-template/config"
	"go.uber.org/zap"
	"gopkg.in/Shopify/sarama.v1"
	"time"
)

type KafkaTopicProducer struct {
	BrokerURL string
	TopicName string
	TopicChan chan *sarama.ProducerMessage
}

// map[Topic_Name] -> Producer
var KafkaTopicProducers = map[string]*KafkaTopicProducer{}

func StartProducers() {
	kafka_broker := config.Config.KafkaBrokerURL
	producer_topics := config.Config.ProducerTopics

	zap.S().Info("Start Producer: kafka_broker=", kafka_broker, " producer_topics=", producer_topics)

	for _, t := range producer_topics {
		// Todo: Register schema

		KafkaTopicProducers[t] = &KafkaTopicProducer{
			kafka_broker,
			t,
			make(chan *sarama.ProducerMessage),
		}

		go KafkaTopicProducers[t].produceTopic()
	}
}

func (k *KafkaTopicProducer) produceTopic() {
	sarama_config := sarama.NewConfig()
	sarama_config.Producer.Partitioner = sarama.NewRandomPartitioner
	sarama_config.Producer.RequiredAcks = sarama.WaitForAll
	sarama_config.Producer.Return.Successes = true

	producer, err := getProducer(k, sarama_config)
	if err != nil {
		zap.S().Info("KAFKA PRODUCER NEWSYNCPRODUCER: Finally Connection cannot be established")
	} else {
		zap.S().Info("KAFKA PRODUCER NEWSYNCPRODUCER: Finally Connection established")
	}
	defer func() {
		if err := producer.Close(); err != nil {
			zap.S().Panic("KAFKA PRODUCER CLOSE PANIC: ", err.Error())
		}
	}()

	zap.S().Debug("Producer ", k.TopicName, ": Started producing")
	for {
		topic_msg := <-k.TopicChan

		partition, offset, err := producer.SendMessage(topic_msg)
		if err != nil {
			zap.S().Warn("Producer ", k.TopicName, ": Err sending message=", err.Error())
		}

		zap.S().Debug("Producer ", k.TopicName, ": Producing message partition=", partition, " offset=", offset)
	}
}

func getProducer(k *KafkaTopicProducer, sarama_config *sarama.Config) (sarama.SyncProducer, error) {
	var producer sarama.SyncProducer
	operation := func() error {
		pro, err := sarama.NewSyncProducer([]string{k.BrokerURL}, sarama_config)
		if err != nil {
			zap.S().Info("KAFKA PRODUCER NEWSYNCPRODUCER PANIC: ", err.Error())
		} else {
			producer = pro
		}
		return err
	}
	neb := backoff.NewExponentialBackOff()
	neb.MaxElapsedTime = time.Minute
	err := backoff.Retry(operation, neb)
	return producer, err
}
