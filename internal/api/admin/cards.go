package admin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pedrofaria/eventcard/internal/service"
)

type CardIdUri struct {
	Id uint32 `uri:"id" binding:"required,numeric,min=1,max=9999999999"`
}

type NewCardData struct {
	ExternalID uint32 `json:"external_id" binding:"required,numeric,min=1,max=9999999999"`
	Name       string `json:"name" binding:"required,min=4,max=80"`
	Enabled    bool   `json:"enabled"`
}

type CardsHandler struct {
	cardService *service.Card
}

func Register(r *gin.Engine, cardService *service.Card) {
	h := &CardsHandler{
		cardService: cardService,
	}

	grp := r.Group("/admin")
	{
		grp.POST("/card", h.CreateCard)
		grp.GET("/card/:id", h.GetCard)
		grp.GET("/card/:id/balance", h.GetCardBalance)
	}
}

func (h *CardsHandler) GetCard(c *gin.Context) {
	ctx := c.Request.Context()

	var uri CardIdUri
	if err := c.BindUri(&uri); err != nil {
		return
	}

	card, err := h.cardService.GetCardFull(ctx, int32(uri.Id))
	if err != nil {
		if errors.Is(err, service.ErrCardNotFound) {
			c.JSON(http.StatusNotFound, nil)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, card)
}

func (h *CardsHandler) GetCardBalance(c *gin.Context) {
	ctx := c.Request.Context()

	var uri CardIdUri
	if err := c.BindUri(&uri); err != nil {
		return
	}

	balance, err := h.cardService.GetCardBalance(ctx, int32(uri.Id))
	if err != nil {
		if errors.Is(err, service.ErrCardNotFound) {
			c.JSON(http.StatusNotFound, nil)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balance)
}

func (h *CardsHandler) CreateCard(c *gin.Context) {
	ctx := c.Request.Context()

	var payload NewCardData

	if err := c.Bind(&payload); err != nil {
		return
	}

	card, err := h.cardService.CreateCard(ctx, payload.ExternalID, payload.Name, payload.Enabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, card)
}
