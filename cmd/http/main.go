package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/pedrofaria/eventcard/internal/api/admin"
	"github.com/pedrofaria/eventcard/internal/bundles/deposit"
	depositApiAdmin "github.com/pedrofaria/eventcard/internal/bundles/deposit/api/admin"
	"github.com/pedrofaria/eventcard/internal/repository"
	"github.com/pedrofaria/eventcard/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "user=eventcard dbname=eventcard password=eventcard host=db sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	cardsRepo := repository.NewCards(db)
	cardService := service.NewCard(cardsRepo)

	r := gin.Default()

	admin.Register(r, cardService)

	depositService := deposit.Init(db, cardService)
	depositApiAdmin.Register(r, depositService)

	if err := r.Run(); err != nil {
		panic(err)
	}
}
