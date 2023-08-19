package broker

import (
	"encoding/json"
	"order/internal/service"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type PaymentStatusEvent struct {
	Data struct {
		OrderID int64  `json:"order_id"`
		Status  string `json:"status"`
		Reason  string `json:"reason"`
	} `json:"data"`
}

type PaymentStatusHandler struct {
	services *service.Services
}

func BuildPaymentStatusHandler(services *service.Services) PaymentStatusHandler {
	return PaymentStatusHandler{services: services}
}

func (gch PaymentStatusHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		gce := PaymentStatusEvent{}
		err := json.Unmarshal(msg.Value, &gce)
		if err != nil {
			logrus.Errorf("Event hasn't been handled, error =  %s", err.Error())
			session.MarkMessage(msg, "")
			continue
		}

		// TODO: update order status

		session.MarkMessage(msg, "")
	}

	return nil
}

func (PaymentStatusHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (PaymentStatusHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
