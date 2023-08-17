package handler

import (
	"user/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/user")
	{
		api.POST("/", h.create)
		api.GET("/:id", h.getById)
		api.GET("/", h.getAll)
		api.DELETE("/:id", h.delete)
		api.PUT("/:id", h.update)

	}

	return router
}
