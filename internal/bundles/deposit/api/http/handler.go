package admin

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mainAdmin "github.com/pedrofaria/eventcard/internal/api/admin"
	"github.com/pedrofaria/eventcard/internal/bundles/deposit"
	"github.com/pedrofaria/eventcard/internal/bundles/deposit/api/http/dto"
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
		grp.GET("/card/:id/deposit", h.List)
		grp.POST("/card/:id/deposit", h.CreateCardDeposit)
		grp.POST("/card/:id/deposit/:depositId/pay", h.TooglePay)
		grp.POST("/card/:id/deposit/:depositId/cancel", h.Cancel)
	}
}

type CardIdDepositIdUri struct {
	mainAdmin.CardIdUri

	DepositId string `uri:"depositId" binding:"required,uuid"`
}

func (d CardIdDepositIdUri) DepositUUID() uuid.UUID {
	return uuid.MustParse(d.DepositId)
}

type NewDepositData struct {
	Amount float32 `json:"amount" binding:"required,gte=0.011,max=99999999.99"`
	Paid   bool    `json:"paid"`
}

func (h *handler) List(c *gin.Context) {
	ctx := c.Request.Context()

	var uri mainAdmin.CardIdUri
	if err := c.BindUri(&uri); err != nil {
		log.Printf("Error: %v", err)
		return
	}

	dd, err := h.srv.List(ctx, uri.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]dto.Deposit, len(dd))
	for i, d := range dd {
		dtos[i] = dto.FromDepositModel(&d)
	}

	c.JSON(http.StatusOK, gin.H{"data": dtos})
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

	d, err := h.srv.CreateDeposit(ctx, uint32(uri.Id), payload.Amount, payload.Paid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.FromDepositModel(d))
}

func (h *handler) TooglePay(c *gin.Context) {
	ctx := c.Request.Context()

	var uri CardIdDepositIdUri
	if err := c.BindUri(&uri); err != nil {
		log.Printf("Error: %v", err)
		return
	}

	d, err := h.srv.TogglePaidStatus(ctx, uint32(uri.Id), uri.DepositUUID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.FromDepositModel(d))
}

func (h *handler) Cancel(c *gin.Context) {
	ctx := c.Request.Context()

	var uri CardIdDepositIdUri
	if err := c.BindUri(&uri); err != nil {
		log.Printf("Error: %v", err)
		return
	}

	err := h.srv.Cancel(ctx, uint32(uri.Id), uri.DepositUUID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, nil)
}
