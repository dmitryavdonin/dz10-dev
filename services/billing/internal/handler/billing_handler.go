package handler

import (
	"net/http"
	"strconv"
	"time"

	"billing/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// create an account for user with user_id
func (h *Handler) createAccount(c *gin.Context) {
	logrus.Printf("createAccount(): BEGIN")

	var input model.NewAccount
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logrus.Errorf("createAccount(): Cannot parse input, error = %s", err.Error())
		return
	}

	item := model.BillingAccount{
		UserId:     input.UserId,
		Balance:    0,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	id, err := h.services.Account.Create(item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logrus.Errorf("createAccount(): Cannot create an account for user with user_id = %d, error = %s", input.UserId, err.Error())
		return
	}

	item.ID = id

	c.JSON(http.StatusOK, item)

	logrus.Printf("createAccount(): END account '%d' is created for user_id = %d", id, input.UserId)
}

// get account for user_id
func (h *Handler) getById(c *gin.Context) {
	logrus.Printf("getById(): BEGIN")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logrus.Errorf("getById(): Cannot parse id, error = %s", err.Error())
		return
	}

	item, err := h.services.Account.GetById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logrus.Errorf("getById(): Cannot get created account for user_id = %d, error = %s", id, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)

	logrus.Printf("getById(): END account for user_id = %d, balance = %d", id, item.Balance)
}

func (h *Handler) getTransactionsForUserId(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logrus.Errorf("getById(): Cannot parse id, error = %s", err.Error())
		return
	}

	var page = c.DefaultQuery("page", "1")
	var limit = c.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	items, err := h.services.Transaction.GetAll(id, intLimit, offset)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "results": len(items), "data": items})
}

func (h *Handler) getAll(c *gin.Context) {

	var page = c.DefaultQuery("page", "1")
	var limit = c.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var items []model.BillingAccount
	items, err := h.services.Account.GetAll(intLimit, offset)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "results": len(items), "data": items})
}

func (h *Handler) delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.services.Account.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Deleted")
}

// put the money on the user's deposit
func (h *Handler) deposit(c *gin.Context) {
	var input model.Deposit
	logrus.Printf("deposit(): BEGIN")
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logrus.Printf("deposit(): put the money on the user's account for user_id = %d, amount = %d", input.UserId, input.Amount)

	item, err := h.services.Account.GetById(input.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logrus.Errorf("deposit(): Cannot get account for user_id = %d, error = %s", input.UserId, err.Error())
		return
	}

	tr := model.BillingTransaction{
		OrderId:    item.ID,
		UserId:     input.UserId,
		Operation:  "deposit",
		Amount:     input.Amount,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	item.Balance += input.Amount
	item.ModifiedAt = time.Now()

	if err := h.services.Account.Update(input.UserId, item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logrus.Errorf("deposit(): Cannot update account for user_id = %d, error = %s", input.UserId, err.Error())

		tr.Status = "failed"
		tr.Reason = err.Error()

		tr_id, err := h.services.Transaction.Create(tr)
		if err != nil {
			logrus.Errorf("deposit(): Cannot create transaction for user_id = %d, error = %s", input.UserId, err.Error())
		} else {
			logrus.Printf("deposit(): Transaction created id = %d user_id = %d, status = %s, reason = %s",
				tr_id, input.UserId, tr.Status, tr.Reason)
		}

		return
	}

	tr.Status = "success"

	tr_id, err := h.services.Transaction.Create(tr)
	if err != nil {
		logrus.Errorf("deposit(): Cannot create transaction for user_id = %d, error = %s", input.UserId, err.Error())
	} else {
		logrus.Printf("deposit(): Transaction created id = %d user_id = %d, status = %s, reason = %s",
			tr_id, input.UserId, tr.Status, tr.Reason)
	}

	logrus.Printf("deposit(): END user_id = %d, new balance = %d", item.UserId, item.Balance)

	c.JSON(http.StatusOK, item)
}

// get money from user's account (withdrawal)
func (h *Handler) withdrawal(c *gin.Context) {
	var input model.Withdrawal
	logrus.Printf("withdrawal(): BEGIN")
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logrus.Printf("withdrawal(): get the money from the user's account for user_id = %d, amount = %d", input.UserId, input.Amount)

	item, err := h.services.Account.GetById(input.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logrus.Errorf("withdrawal(): Cannot get account for user_id = %d, error = %s", input.UserId, err.Error())
		return
	}

	tr := model.BillingTransaction{
		OrderId:    item.ID,
		UserId:     input.UserId,
		Operation:  "withdrawal",
		Amount:     input.Amount,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	if item.Balance < input.Amount {
		statusResponse := StatusResponse{
			Status: "failed",
			Reason: "Out of balance",
		}

		c.JSON(http.StatusOK, statusResponse)

		logrus.Printf("withdrawal(): Cannot withraw money from user_id = %d, status = %s, reason = %s",
			input.UserId, statusResponse.Status, statusResponse.Reason)

		tr.Status = statusResponse.Status
		tr.Reason = statusResponse.Reason

		tr_id, err := h.services.Transaction.Create(tr)
		if err != nil {
			logrus.Errorf("withdrawal(): Cannot create transaction for user_id = %d, error = %s", input.UserId, err.Error())
		} else {
			logrus.Printf("withdrawal(): Transaction created id = %d user_id = %d, status = %s, reason = %s",
				tr_id, input.UserId, tr.Status, tr.Reason)
		}

		return
	}

	item.Balance -= input.Amount
	item.ModifiedAt = time.Now()

	if err := h.services.Account.Update(input.UserId, item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logrus.Errorf("withdrawal(): Cannot update account for user_id = %d, error = %s", input.UserId, err.Error())

		tr.Status = "failed"
		tr.Reason = err.Error()

		tr_id, err := h.services.Transaction.Create(tr)
		if err != nil {
			logrus.Errorf("withdrawal(): Cannot create transaction for user_id = %d, error = %s", input.UserId, err.Error())
		} else {
			logrus.Printf("withdrawal(): Transaction created id = %d user_id = %d, status = %s, reason = %s",
				tr_id, input.UserId, tr.Status, tr.Reason)
		}

		return
	}

	tr.Status = "success"

	tr_id, err := h.services.Transaction.Create(tr)
	if err != nil {
		logrus.Errorf("withdrawal(): Cannot create transaction for user_id = %d, error = %s", input.UserId, err.Error())
	} else {
		logrus.Printf("withdrawal(): Transaction created id = %d user_id = %d, status = %s, reason = %s",
			tr_id, input.UserId, tr.Status, tr.Reason)
	}

	c.JSON(http.StatusOK, item)
}
