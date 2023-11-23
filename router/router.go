package router

import (
	"finalProject4/controller"
	"finalProject4/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("/user")
	{
		userRouter.POST("/register", controller.NewUserHandlerImpl().UserRegister)

		userRouter.POST("/login", controller.NewUserHandlerImpl().UserLogin)

		userRouter.Use(middleware.Authentication())
		userRouter.PATCH("/topup", controller.NewUserHandlerImpl().BalanceUpdate)

		userRouter.DELETE("/:userID", middleware.UserAuthorization(), controller.NewUserHandlerImpl().UserDelete)
	}

	categoryRouter := r.Group("/category")
	{
		categoryRouter.Use(middleware.Authentication())
		categoryRouter.Use(middleware.AdminAuthMiddleware())
		categoryRouter.POST("/create", controller.NewCategoryHandlerImpl().CategoryCreate)
		categoryRouter.GET("/get", controller.NewCategoryHandlerImpl().CategoryGet)
		categoryRouter.PATCH("/update/:categoryID", controller.NewCategoryHandlerImpl().CategoryUpdate)
		categoryRouter.DELETE("/delete/:categoryID", controller.NewCategoryHandlerImpl().CategoryDelete)
	}

	productRouter := r.Group("/product")
	{
		productRouter.Use(middleware.Authentication())
		productRouter.Use(middleware.AdminAuthMiddleware())
		fmt.Println("AdminAuthMiddleware is executing")
		productRouter.POST("/create", middleware.ProductAuthentication(), controller.NewProductHandlerImpl().ProductCreate)
		productRouter.GET("/get", controller.NewProductHandlerImpl().ProductGet)
		productRouter.PUT("/update/:productID", middleware.ProductAuthentication(), controller.NewProductHandlerImpl().ProductUpdate)
		productRouter.DELETE("/delete/:productID", controller.NewProductHandlerImpl().ProductDelete)
	}

	transactionRouter := r.Group("/transactions")
	{
		transactionRouter.Use(middleware.Authentication())
		transactionRouter.POST("/create", middleware.TransactionAuthentication(), controller.NewTransactionHandlerImpl().TransactionCreate)
		transactionRouter.GET("/my-transactions", controller.NewTransactionHandlerImpl().TransactionGetMy)
		transactionRouter.GET("/user-transactions", middleware.AdminAuthMiddleware(), controller.NewTransactionHandlerImpl().TransactionGetAll)
	}

	return r
}
