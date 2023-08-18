package handler

import (
	"order/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services: services}
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
