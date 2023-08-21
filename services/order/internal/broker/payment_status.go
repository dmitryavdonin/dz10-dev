package broker

import (
	"encoding/json"
	"order/internal/service"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type PaymentStatusEvent struct {
	Data struct {
		OrderId int    `json:"order_id"`
		UserId  int    `json:"user_id"`
		Status  string `json:"status"`
		Reason  string `json:"reason"`
	} `json:"data"`
}

type PaymentStatusHandler struct {
	services           *service.Services
	paymentStatusTopic string
}

func BuildPaymentStatusHandler(services *service.Services, paymentStatusTopic string) PaymentStatusHandler {
	return PaymentStatusHandler{
		services:           services,
		paymentStatusTopic: paymentStatusTopic,
	}
}

func (psh PaymentStatusHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	logrus.Printf("ConsumeClaim(): BEGIN consuming messages from topic = %s", psh.paymentStatusTopic)

	for msg := range claim.Messages() {
		pse := PaymentStatusEvent{}
		err := json.Unmarshal(msg.Value, &pse)
		if err != nil {
			logrus.Errorf("ConsumeClaim(): Event hasn't been handled, error =  %s", err.Error())
			session.MarkMessage(msg, "")
			continue
		}

		logrus.Printf("ConsumeClaim(): Message received: order_id = %d, user_id = %d, status = %s, reason = %s",
			pse.Data.OrderId, pse.Data.UserId, pse.Data.Status, pse.Data.Reason)

		// update order status with the results of payment status
		logrus.Printf("ConsumeClaim(): Try to update order status order_id = %d, status = %s, reason = %s",
			pse.Data.OrderId, pse.Data.Status, pse.Data.Reason)

		order, err := psh.services.Order.GetById(pse.Data.OrderId)
		if err != nil {
			logrus.Errorf("ConsumeClaim(): Cannot get order by order_id = %d, error =  %s", pse.Data.OrderId, err.Error())
			session.MarkMessage(msg, "")
			continue
		}

		order.Status = pse.Data.Status
		order.Reason = pse.Data.Reason
		order.ModifiedAt = time.Now()

		if err := psh.services.Order.Update(order.ID, order); err != nil {
			logrus.Errorf("ConsumeClaim(): Cannot get order by order_id = %d, error =  %s", pse.Data.OrderId, err.Error())
			session.MarkMessage(msg, "")
			continue
		}

		logrus.Printf("ConsumeClaim(): END Order updated order_id = %d, status = %s, reason = %s",
			pse.Data.OrderId, pse.Data.Status, pse.Data.Reason)

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
