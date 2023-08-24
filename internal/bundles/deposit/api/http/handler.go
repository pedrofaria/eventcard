package admin

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	mainAdmin "github.com/pedrofaria/eventcard/internal/api/admin"
	"github.com/pedrofaria/eventcard/internal/bundles/deposit"
)

type handler struct {
	srv *deposit.Service
}

func Register(r *gin.Engine, srv *deposit.Service) {
	h := handler{
		srv: srv,
	}

	grp := r.Group("/admin")
	{
		grp.POST("/card/:id/deposit", h.CreateCardDeposit)
	}
}

type NewDepositData struct {
	Amount float32 `json:"amount" binding:"required,gte=0.011,max=99999999.99"`
	Paid   bool    `json:"paid"`
}

func (h *handler) CreateCardDeposit(c *gin.Context) {
	ctx := c.Request.Context()

	var uri mainAdmin.CardIdUri
	if err := c.BindUri(&uri); err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var payload NewDepositData

	if err := c.Bind(&payload); err != nil {
		log.Printf("Error: %v", err)
		return
	}

	deposit, err := h.srv.CreateDeposit(ctx, uint32(uri.Id), payload.Amount, payload.Paid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deposit)
}
