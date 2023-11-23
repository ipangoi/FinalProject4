package controller

import "github.com/gin-gonic/gin"

type TransactionHandler interface {
	TransactionCreate(*gin.Context)
	TransactionGetMy(*gin.Context)
	TransactionGetAll(*gin.Context)
}
