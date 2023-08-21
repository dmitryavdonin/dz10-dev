package broker

import (
	"billing/internal/model"
	"billing/internal/service"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type OrderCreatedEvent struct {
	Data struct {
		OrderId int `json:"order_id"`
		UserId  int `json:"user_id"`
		Price   int `json:"price"`
	} `json:"data"`
}

type OrderCreatedHandler struct {
	services          *service.Services
	kafkaProducer     *KafkaProducer
	orderCreatedTopic string
}

func BuildOrderCreatedHandler(services *service.Services, kafkaProducer *KafkaProducer, orderCreatedTopic string) OrderCreatedHandler {
	return OrderCreatedHandler{
		services:          services,
		kafkaProducer:     kafkaProducer,
		orderCreatedTopic: orderCreatedTopic,
	}
}

func (och OrderCreatedHandler) sendPaymentStatusEvent(tr model.BillingTransaction) error {
	logrus.Print("sendPaymentStatusEvent(): BEGIN")
	msg := model.PaymentStatusMsg{Data: model.PaymentStatus{
		OrderId: tr.OrderId,
		UserId:  tr.UserId,
		Status:  tr.Status,
		Reason:  tr.Reason,
	}}

	msgStr, err := json.Marshal(msg)
	if err != nil {
		logrus.Errorf("sendPaymentStatusEvent(): Cannot marshal the message, error = %s", err.Error())
		return err
	}
	producerMsg := &sarama.ProducerMessage{Topic: och.kafkaProducer.PaymentStatusTopic, Value: sarama.StringEncoder(msgStr)}
	_, _, err = och.kafkaProducer.Producer.SendMessage(producerMsg)
	if err != nil {
		logrus.Errorf("sendPaymentStatusEvent(): Cannot send the message, error = %s", err.Error())
		return err
	} else {
		logrus.Print("sendPaymentStatusEvent(): END")
		return nil
	}
}

func (och OrderCreatedHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	logrus.Printf("ConsumeClaim(): BEGIN consuming messages from topic = %s", och.orderCreatedTopic)
	for msg := range claim.Messages() {
		oce := OrderCreatedEvent{}
		err := json.Unmarshal(msg.Value, &oce)
		if err != nil {
			logrus.Errorf("ConsumeClaim(): Event hasn't been handled, error =  %s", err.Error())
			session.MarkMessage(msg, "")
			continue
		}

		logrus.Printf("ConsumeClaim(): Message received: order_id = %d, user_id = %d, price = %d",
			oce.Data.OrderId, oce.Data.UserId, oce.Data.Price)

		// get user balance
		logrus.Printf("ConsumeClaim(): Try to get user balance for user_id = %d", oce.Data.UserId)
		account, err := och.services.Account.GetById(oce.Data.UserId)
		if err != nil {
			logrus.Errorf("ConsumeClaim(): Cannot get user balance for user_id = %d, error =  %s", oce.Data.UserId, err.Error())
			session.MarkMessage(msg, "")
			continue
		}

		// if there are enough money, then update account with new balance, and create a transaction withdrawal operation
		if account.Balance >= oce.Data.Price {
			account.Balance -= oce.Data.Price

			logrus.Printf("ConsumeClaim(): Try to update user balance for user_id = %d, order_id = %d, new balance = %d",
				oce.Data.UserId, oce.Data.OrderId, account.Balance)

			if err := och.services.Account.Update(account.UserId, account); err != nil {
				logrus.Errorf("ConsumeClaim(): Cannot update user balance for user_id = %d, new balance = %d, error =  %s",
					oce.Data.UserId, account.Balance, err.Error())
				session.MarkMessage(msg, "")
				continue
			}

			tr := model.BillingTransaction{
				UserId:    oce.Data.UserId,
				OrderId:   oce.Data.OrderId,
				Amount:    oce.Data.Price,
				Operation: "withdrawal",
				Status:    "success",
			}

			logrus.Printf("ConsumeClaim(): Try to create transaction for user_id = %d, order_id = %d, operation = %s, amount = %d",
				tr.UserId, tr.OrderId, tr.Operation, tr.Amount)
			if _, err := och.services.Transaction.Create(tr); err != nil {
				logrus.Errorf("ConsumeClaim(): Cannot create transaction for user_id = %d, order_id = %d, operation = %s, amount = %d, error =  %s",
					tr.UserId, tr.OrderId, tr.Operation, tr.Amount, err.Error())
			}
			session.MarkMessage(msg, "")

			logrus.Printf("ConsumeClaim(): Try to send payment_status_event user_id = %d, order_id = %d, operation = %s, amount = %d",
				tr.UserId, tr.OrderId, tr.Operation, tr.Amount)

			if err := och.sendPaymentStatusEvent(tr); err != nil {
				logrus.Errorf("ConsumeClaim(): Cannot send payment_status_event for user_id = %d, order_id = %d, operation = %s, amount = %d, error =  %s",
					tr.UserId, tr.OrderId, tr.Operation, tr.Amount, err.Error())
			}

			continue
		} else {
			// balance is not enough for order
			logrus.Printf("ConsumeClaim(): User balance = %d for user_id = %d, order_id = %d, is not enough for price = %d",
				account.Balance, oce.Data.UserId, oce.Data.OrderId, oce.Data.Price)

			tr := model.BillingTransaction{
				UserId:    oce.Data.UserId,
				OrderId:   oce.Data.OrderId,
				Amount:    oce.Data.Price,
				Operation: "withdrawal",
				Status:    "failed",
				Reason:    "Out of balance",
			}

			logrus.Printf("ConsumeClaim(): Try to create transaction for user_id = %d, order_id = %d, operation = %s, amount = %d",
				tr.UserId, tr.OrderId, tr.Operation, tr.Amount)
			if _, err := och.services.Transaction.Create(tr); err != nil {
				logrus.Errorf("ConsumeClaim(): Cannot create transaction for user_id = %d, order_id = %d, operation = %s, amount = %d, error =  %s",
					tr.UserId, tr.OrderId, tr.Operation, tr.Amount, err.Error())
			}
			session.MarkMessage(msg, "")

			logrus.Printf("ConsumeClaim(): Try to send payment_status_event user_id = %d, order_id = %d, operation = %s, amount = %d",
				tr.UserId, tr.OrderId, tr.Operation, tr.Amount)

			if err := och.sendPaymentStatusEvent(tr); err != nil {
				logrus.Errorf("ConsumeClaim(): Cannot send payment_status_event for user_id = %d, order_id = %d, operation = %s, amount = %d, error =  %s",
					tr.UserId, tr.OrderId, tr.Operation, tr.Amount, err.Error())
			}
			continue
		}
	}

	logrus.Printf("ConsumeClaim(): END")

	return nil
}

func (OrderCreatedHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (OrderCreatedHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
