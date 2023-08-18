package handler

import (
	"billing/internal/service"

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
		api.POST("/account", h.createAccount)
		api.GET("/account/:id", h.getById)
		api.GET("/account", h.getAll)
		api.DELETE("/account/:id", h.delete)
		api.POST("/deposit", h.deposit)
		api.POST("/withdrawal", h.withdrawal)
		api.GET("/transaction/:id", h.getTransactionsForUserId)

	}

	return router
}
