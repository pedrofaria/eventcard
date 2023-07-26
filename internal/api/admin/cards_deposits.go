package admin

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NewDepositData struct {
	Amount float32 `json:"amount" binding:"required,gte=0.011,max=99999999.99"`
	Paid   bool    `json:"paid"`
}

func (h *CardsHandler) CreateCardDeposit(c *gin.Context) {
	ctx := c.Request.Context()

	var uri CardIdUri
	if err := c.BindUri(&uri); err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var payload NewDepositData

	if err := c.Bind(&payload); err != nil {
		log.Printf("Error: %v", err)
		return
	}

	deposit, err := h.cardService.CreateDeposit(ctx, uint32(uri.Id), payload.Amount, payload.Paid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deposit)
}
