package broker

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

func RunConsumers(ctx context.Context, handlers map[string]sarama.ConsumerGroupHandler, kafka_host string, kafka_port string, consume_topic string) {
	kafkaConsumerGroups := initAllConsumerGroups(kafka_host, kafka_port, consume_topic)

	for topic, group := range kafkaConsumerGroups {
		go func(topic string, group *sarama.ConsumerGroup) {
			defer func() {
				if r := recover(); r != nil {
					logrus.Fatalf("Cannot recover consumer group: %s", r)
				}
			}()

			for {
				err := (*group).Consume(ctx, []string{topic}, handlers[topic])
				if err != nil {
					logrus.Errorf("Consumer group error =  %s", err.Error())
				}
			}
		}(topic, group)
	}
}

func initAllConsumerGroups(kafka_host string, kafka_port string, consume_topic string) map[string]*sarama.ConsumerGroup {
	topic := consume_topic
	return map[string]*sarama.ConsumerGroup{
		topic: initGroup(kafka_host, kafka_port, consume_topic),
	}
}

func initGroup(kafka_host string, kafka_port string, consume_topic string) *sarama.ConsumerGroup {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_3_0_0
	cfg.Consumer.Return.Errors = true

	logrus.Printf("Try to connect to consumer group in kafka host = %s, port = %s, topic = %s",
		kafka_host, kafka_port, consume_topic)

	group, err := sarama.NewConsumerGroup([]string{kafka_host + ":" + kafka_port}, consume_topic, cfg)
	if err != nil {
		logrus.Errorf("Message hasn't been marshaled, error =  %s", err.Error())
		return nil
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Fatalf("Cannot recover consumer group: %s", r)
			}
		}()

		for err := range group.Errors() {
			logrus.Errorf("Consumer group error =  %s", err.Error())
		}
	}()

	return &group
}
