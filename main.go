package main

import (
	"finalProject4/database"
	"finalProject4/entity"
	"finalProject4/router"
	"fmt"
	"log"

	"gorm.io/gorm"
)

func main() {
	database.StartDB()
	var db = database.GetDB()
	seedAdminData(db)
	r := router.StartApp()
	r.Run(":8080")
}

func seedAdminData(db *gorm.DB) {
	var adminCount int64
	var User []entity.User
	db.Model(&User).Where("role = ?", "admin").Count(&adminCount)

	if adminCount == 0 {
		admin := entity.User{
			Full_Name: "admin",
			Email:     "admin@gmail.com",
			Password:  "admin123",
			Role:      "admin",
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatal(err)
		}

		fmt.Println("Admin user seeded successfully.")
	} else {
		fmt.Println("Admin user already exists.")
	}
}
