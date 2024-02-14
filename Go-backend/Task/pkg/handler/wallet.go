package handler

import (
	"errors"
	"github.com/geejjoo/task"
	"github.com/geejjoo/task/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) createWallet(c *gin.Context) {
	id, balance, err := h.services.CreateWallet()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"balance": balance,
	})

}

func (h *Handler) getWallet(c *gin.Context) {
	idWallet := c.Param("id")
	id, balance, err := h.services.GetWallet(idWallet)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"balance": balance,
	})
}

func (h *Handler) getHistory(c *gin.Context) {
	idWallet := c.Param("id")
	history, err := h.services.GetHistory(idWallet)
	if errors.Is(err, service.WalletNotFoundError) {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, history)
}

func (h *Handler) updateWallet(c *gin.Context) {
	var input task.SendRequest
	idWallet := c.Param("id")
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest,
			"Bad request")
		return
	}

	updateWallet := task.UpdateWallet{
		FromID: idWallet,
		ToID:   input.ToID,
		Amount: input.Amount,
	}

	if updateWallet.FromID == updateWallet.ToID {
		newErrorResponse(c, http.StatusNotFound, "Send wallet and target wallet are the same")
		return
	}
	if updateWallet.Amount <= 0 {
		newErrorResponse(c, http.StatusBadRequest, "Amount should be positive")
		return
	}

	err := h.services.UpdateWallet(&updateWallet)

	switch {
	case errors.Is(err, service.FromIdError):
		newErrorResponse(c, http.StatusNotFound, "Incorrect wallet ID")
		return
	case errors.Is(err, service.ToIdError):
		newErrorResponse(c, http.StatusBadRequest, "Incorrect target ID")
		return
	case errors.Is(err, service.BalanceError):
		newErrorResponse(c, http.StatusBadRequest, "Not enough money")
		return
	}

	c.Status(http.StatusOK)
}
