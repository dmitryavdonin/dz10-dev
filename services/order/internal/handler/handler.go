package handler

import (
	"order/internal/broker"
	"order/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services      *service.Services
	kafkaProducer *broker.KafkaProducer
}

func NewHandler(services *service.Services, kafkaProducer *broker.KafkaProducer) *Handler {
	return &Handler{
		services:      services,
		kafkaProducer: kafkaProducer,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/")
	{
		api.POST("/", h.createOrder)
		api.GET("/:id", h.getOrderById)
		api.GET("/", h.getAllOrders)
		api.DELETE("/:id", h.deleteOrder)
	}

	return router
}
