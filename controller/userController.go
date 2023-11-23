package controller

import "github.com/gin-gonic/gin"

type UserHandler interface {
	UserRegister(*gin.Context)
	UserLogin(*gin.Context)
	BalanceUpdate(*gin.Context)
	UserDelete(*gin.Context)
}
