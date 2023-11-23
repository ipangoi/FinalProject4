package controller

import "github.com/gin-gonic/gin"

type ProductHandler interface {
	ProductCreate(*gin.Context)
	ProductGet(*gin.Context)
	ProductUpdate(*gin.Context)
	ProductDelete(*gin.Context)
}
