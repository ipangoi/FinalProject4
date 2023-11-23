package controller

import (
	"encoding/json"
	"finalProject4/database"
	"finalProject4/entity"
	"finalProject4/helper"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandlerImpl struct{}

func NewProductHandlerImpl() ProductHandler {
	return &ProductHandlerImpl{}
}

func (s *ProductHandlerImpl) ProductCreate(c *gin.Context) {
	var db = database.GetDB()
	categoryData := c.MustGet("categoryData").(map[string]interface{})
	contentType := helper.GetContentType(c)

	categoryID := uint(categoryData["id"].(uint))
	Product := entity.Product{}

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
		if err := json.Unmarshal(rawJSON, &Product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid JSON payload for Comment",
			})
			return
		}
	} else {
		c.ShouldBind(&Product)
	}

	if err := Product.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	Product.CategoryID = categoryID

	err = db.Debug().Create(&Product).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          Product.ID,
		"title":       Product.Title,
		"price":       Product.Price,
		"stock":       Product.Stock,
		"category_id": categoryID,
		"created_at":  Product.CreatedAt,
	})
}

func (s *ProductHandlerImpl) ProductGet(c *gin.Context) {
	var db = database.GetDB()
	contentType := helper.GetContentType(c)

	var Product []entity.Product

	if contentType == appJSON {
		c.ShouldBindJSON(&Product)
	} else {
		c.ShouldBind(&Product)
	}

	err := db.Find(&Product).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	var response []gin.H
	for _, product := range Product {
		response = append(response, gin.H{
			"id":          product.ID,
			"title":       product.Title,
			"price":       product.Price,
			"stock":       product.Stock,
			"category_id": product.CategoryID,
			"created_at":  product.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (s *ProductHandlerImpl) ProductUpdate(c *gin.Context) {
	var db = database.GetDB()
	categoryData := c.MustGet("categoryData").(map[string]interface{})
	contentType := helper.GetContentType(c)
	_, _ = db, contentType

	categoryID := uint(categoryData["id"].(uint))

	Product := entity.Product{}

	productID, _ := strconv.Atoi(c.Param("productID"))

	if contentType == appJSON {
		c.ShouldBindJSON(&Product)
	} else {
		c.ShouldBind(&Product)
	}

	Product.ID = uint(productID)

	err := db.Model(&Product).Where("id = ?", productID).Updates(
		entity.Product{
			Title:      Product.Title,
			Price:      Product.Price,
			Stock:      Product.Stock,
			CategoryID: categoryID}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          Product.ID,
		"title":       Product.Title,
		"price":       Product.Price,
		"stock":       Product.Stock,
		"category_id": Product.CategoryID,
		"created_at":  Product.CreatedAt,
		"updated_at":  Product.UpdatedAt,
	})
}

func (s *ProductHandlerImpl) ProductDelete(c *gin.Context) {
	var db = database.GetDB()
	contentType := helper.GetContentType(c)

	Product := entity.Product{}

	productID, _ := strconv.Atoi(c.Param("productID"))

	if contentType == appJSON {
		c.ShouldBindJSON(&Product)
	} else {
		c.ShouldBind(&Product)
	}

	Product.ID = uint(productID)

	err := db.Model(&Product).Where("id = ?", productID).Delete(&Product).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product has been successfully deleted",
	})
}
