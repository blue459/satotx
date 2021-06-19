package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"log"
	"satotx/controller"
	"time"

	"github.com/gin-gonic/gin"
	ginadapter "github.com/linthan/scf-go-api-proxy/gin"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tencentyun/scf-go-lib/events"
)

var ginLambda *ginadapter.GinLambda

func init() {
	log.Println("Gin cold start")
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type,Origin,Content-Encoding"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
	router.GET("/", controller.Satotx)
	router.POST("/utxo/:txid/:index", controller.SignUtxo)
	router.POST("/utxo-spend-by/:txid/:index/:byTxid", controller.SignUtxoSpendBy)
	router.POST("/utxo-spend-by-utxo/:txid/:index/:byTxid/:byTxindex", controller.SignUtxoSpendByUtxo)

	ginLambda = ginadapter.New(router)
}

func Handler(ctx context.Context, req events.APIGatewayRequest) (events.APIGatewayResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	cloudfunction.Start(Handler)
}
