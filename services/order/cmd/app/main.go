package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"order/internal/broker"
	"order/internal/config"
	"order/internal/handler"
	"order/internal/order"
	"order/internal/repository"
	"order/internal/service"

	"github.com/IBM/sarama"
	"github.com/dmitryavdonin/gtools/migrations"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	cfg, err := config.InitConfig("")
	if err != nil {
		panic(fmt.Sprintf("error initializing config %s", err))
	}

	//db migrations
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.DBname)

	migrate, err := migrations.NewMigrations(dsn, "file://migrations")
	if err != nil {
		logrus.Fatalf("migrations error: %s", err.Error())
	}

	err = migrate.Up()
	if err != nil {
		logrus.Fatalf("migrations error: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(dsn)

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	// create repository
	repos := repository.NewRepository(db)
	//create servicec
	services := service.NewServices(repos)
	// create kafaka producer
	producer := broker.InitKafkaProducer(cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.OrderCreatedTopic)
	// create hanlers
	handlers := handler.NewHandler(services, producer)
	//crate kafka consumer
	broker_handlers := map[string]sarama.ConsumerGroupHandler{
		cfg.Kafka.PaymentStatusTopic: broker.BuildPaymentStatusHandler(services),
	}
	broker.RunConsumers(context.Background(), broker_handlers, cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.PaymentStatusTopic)

	// create server
	srv := new(order.Server)
	go func() {
		if err := srv.Run(cfg.App.Port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}

	}()

	logrus.Printf("Service %s started on port = %d ", cfg.App.ServiceName, cfg.App.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Printf("Service %s shutting down", cfg.App.ServiceName)

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
}
