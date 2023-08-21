package main

import (
	"billing"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"billing/internal/broker"
	"billing/internal/config"
	"billing/internal/handler"
	"billing/internal/repository"
	"billing/internal/service"

	"github.com/IBM/sarama"
	"github.com/dmitryavdonin/gtools/migrations"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	cfg, err := config.InitConfig("")
	if err != nil {
		panic(fmt.Sprintf("main(): error initializing config %s", err))
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
		logrus.Fatalf("main(): migrations error: %s", err.Error())
	}

	err = migrate.Up()
	if err != nil {
		logrus.Fatalf("main(): migrations error: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(dsn)

	if err != nil {
		logrus.Fatalf("main(): failed to initialize db: %s", err.Error())
	}

	// create repository
	repos := repository.NewRepository(db)
	//create services
	services := service.NewServices(repos)
	// create kafaka producer
	producer := broker.InitKafkaProducer(cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.PaymentStatusTopic)
	// create hanlers
	handlers := handler.NewHandler(services, producer)
	//crate kafka consumer
	broker_handlers := map[string]sarama.ConsumerGroupHandler{
		cfg.Kafka.OrderCreatedTopic: broker.BuildOrderCreatedHandler(services, producer, cfg.Kafka.OrderCreatedTopic),
	}
	broker.RunConsumers(context.Background(), broker_handlers, cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.OrderCreatedTopic)

	// create server
	srv := new(billing.Server)
	go func() {
		if err := srv.Run(cfg.App.Port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("main(): error occured while running http server: %s", err.Error())
		}

	}()

	logrus.Printf("main(): Service %s started on port = %d ", cfg.App.ServiceName, cfg.App.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Printf("main(): Service %s shutting down", cfg.App.ServiceName)

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("main(): error occured on server shutting down: %s", err.Error())
	}
}
