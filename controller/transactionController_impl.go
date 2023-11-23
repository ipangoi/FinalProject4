package controller

import (
	"encoding/json"
	"finalProject4/database"
	"finalProject4/entity"
	"finalProject4/helper"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionHandlerImpl struct{}

func NewTransactionHandlerImpl() TransactionHandler {
	return &TransactionHandlerImpl{}
}

func (s *TransactionHandlerImpl) TransactionCreate(c *gin.Context) {
	var db = database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	productData := c.MustGet("productData").(map[string]interface{})
	contentType := helper.GetContentType(c)

	userID := uint(userData["id"].(float64))
	productID := uint(productData["id"].(uint))
	Transaction := entity.TransactionHistory{}

	rawJSON, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Error reading raw JSON data",
		})
		return
	}

	fmt.Println("Raw JSON request:", string(rawJSON))
	fmt.Println("Request Headers:", c.Request.Header)

	if contentType == appJSON {
		if err := json.Unmarshal(rawJSON, &Transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid JSON payload for Comment",
			})
			return
		}
	} else {
		c.ShouldBind(&Transaction)
	}

	if err := Transaction.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	Transaction.UserID = userID
	Transaction.ProductID = productID

	var product entity.Product
	if err := db.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Product not found",
		})
		return
	}

	if product.Stock < Transaction.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Insufficient stock",
		})
		return
	}

	err = db.Debug().Preload("User").Where("user_id = ?", userID).First(&Transaction).Error
	if product.Price*Transaction.Quantity > Transaction.User.Balance {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Insufficient balance",
		})
		return
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&product).Update("stock", gorm.Expr("stock - ?", Transaction.Quantity)).Error; err != nil {
			return err
		}

		if err := tx.Model(&entity.User{}).Where("id = ?", userID).Update("balance", gorm.Expr("balance - ?", product.Price*Transaction.Quantity)).Error; err != nil {
			return err
		}

		if err := tx.Model(&entity.Category{}).Where("id = ?", product.CategoryID).Update("sold_product_amount", gorm.Expr("sold_product_amount + ?", Transaction.Quantity)).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	Transaction.Total_Price = product.Price * Transaction.Quantity

	err = db.Debug().Omit("ID").Create(&Transaction).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	err = db.Debug().Preload("Product").First(&Transaction, Transaction.ID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "You have succesfully purchased the product",
		"transaction_bill": gin.H{
			"id":            Transaction.ID,
			"quantity":      Transaction.Quantity,
			"total_price":   Transaction.Total_Price,
			"product_title": Transaction.Product.Title,
		},
	})
}

func (s *TransactionHandlerImpl) TransactionGetMy(c *gin.Context) {
	var db = database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helper.GetContentType(c)

	userID := uint(userData["id"].(float64))

	var Transaction []entity.TransactionHistory

	rawJSON, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Error reading raw JSON data",
		})
		return
	}

	fmt.Println("Raw JSON request:", string(rawJSON))
	fmt.Println("Request Headers:", c.Request.Header)

	if contentType == appJSON {
		if err := json.Unmarshal(rawJSON, &Transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid JSON payload for Comment",
			})
			return
		}
	} else {
		c.ShouldBind(&Transaction)
	}

	for i := range Transaction {
		Transaction[i].UserID = userID

	}

	err = db.Preload("Product").Where("user_id = ?", userID).Find(&Transaction).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Transaction)
}

func (s *TransactionHandlerImpl) TransactionGetAll(c *gin.Context) {
	var db = database.GetDB()
	contentType := helper.GetContentType(c)

	var Transaction []entity.TransactionHistory

	rawJSON, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Error reading raw JSON data",
		})
		return
	}

	fmt.Println("Raw JSON request:", string(rawJSON))
	fmt.Println("Request Headers:", c.Request.Header)

	if contentType == appJSON {
		if err := json.Unmarshal(rawJSON, &Transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid JSON payload for Comment",
			})
			return
		}
	} else {
		c.ShouldBind(&Transaction)
	}

	err = db.Preload("Product").Preload("User").Find(&Transaction).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Transaction)
}
