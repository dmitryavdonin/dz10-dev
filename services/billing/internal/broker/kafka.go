package broker

import (
	"os"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type KafkaProducer struct {
	Producer           sarama.SyncProducer
	PaymentStatusTopic string
}

func InitKafkaProducer(kafka_host string, kafka_port string, topic string) *KafkaProducer {

	brokerCfg := sarama.NewConfig()
	brokerCfg.Producer.RequiredAcks = sarama.WaitForAll
	brokerCfg.Producer.Return.Successes = true

	logrus.Printf("Try to create producer for kafka host = %s, port = %s, topic = %s", kafka_host, kafka_port, topic)

	producer, err := sarama.NewSyncProducer([]string{kafka_host + ":" + kafka_port}, brokerCfg)
	if err != nil {
		logrus.Fatalf("Kafka fatal error = %s", err.Error())
		os.Exit(1)
	}

	return &KafkaProducer{
		PaymentStatusTopic: topic,
		Producer:           producer,
	}
}
