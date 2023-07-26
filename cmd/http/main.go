package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/pedrofaria/eventcard/internal/api/admin"
	"github.com/pedrofaria/eventcard/internal/repository"
	"github.com/pedrofaria/eventcard/internal/service"

	_ "github.com/lib/pq"
)

type CardIdUri struct {
	Id uint32 `uri:"id" binding:"required,numeric,min=1,max=9999999999"`
}

func main() {
	db, err := sql.Open("postgres", "user=eventcard dbname=eventcard password=eventcard host=db sslmode=disable")
	if err != nil {
		panic(err)
	}

	cardsRepo := repository.NewCards(db)
	cardService := service.NewCard(cardsRepo)
	cardsHandler := admin.NewCardsHandler(cardService)

	r := gin.Default()

	admin.Register(r, cardsHandler)

	if err := r.Run(); err != nil {
		panic(err)
	}
}
