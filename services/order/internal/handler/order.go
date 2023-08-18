package handler

import (
	"net/http"
	"strconv"
	"time"

	"order/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// create order
func (h *Handler) createOrder(c *gin.Context) {
	logrus.Print("createOrder(): BEGIN ")
	var input model.NewOrder
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, StatusResponse{Status: "failed", Reason: err.Error()})
		logrus.Errorf("createOrder(): Cannot parse input, error = %s", err.Error())
		return
	}

	now := time.Now()

	order := model.Order{
		UserId:     input.UserId,
		Price:      input.Price,
		Status:     "pending",
		CreatedAt:  now,
		ModifiedAt: now,
	}

	logrus.Printf("createOrder(): Try to create order record: user_id = %d, price = %d", order.UserId, order.Price)
	id, err := h.services.Order.Create(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, StatusResponse{Status: "failed", Reason: err.Error()})
		logrus.Errorf("createOrder(): Cannot create order record, error = %s", err.Error())
		return
	}

	order.ID = id

	c.JSON(http.StatusOK, order)

	logrus.Printf("createOrder(): Try to create Saga for order_id = %d, user_id = %d, price = %d", order.ID, order.UserId, order.Price)
}

// get order by id
func (h *Handler) getOrderById(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.services.Order.GetById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// get all orders
func (h *Handler) getAllOrders(c *gin.Context) {

	var page = c.DefaultQuery("page", "1")
	var limit = c.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var items []model.Order
	items, err := h.services.Order.GetAll(intLimit, offset)
	if err != nil {
		c.JSON(http.StatusBadGateway,
			StatusResponse{
				Status: "failed",
				Reason: err.Error(),
			})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "results": len(items), "data": items})
}

// delete order by id
func (h *Handler) deleteOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.services.Order.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
