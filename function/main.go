package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func main() {
	lambda.Start(Handler)
}

func init() {
	r := gin.Default()

	r.POST("/api/:customerId/transactions", ProcessTxnFile)
	r.GET("/api/:customerId/transactions", GetCustomerTransactions)
	r.GET("/api/:customerId/transactions/:transactionId", GetTransaction)
	
	ginLambda = ginadapter.New(r)

	if isLocal() {
		r.Run(":8080")
	}
	
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}
