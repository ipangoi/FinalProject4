package entity

import "github.com/asaskevich/govalidator"

type TransactionHistory struct {
	GormModel
	Quantity    int `gorm:"not null" json:"quantity" valid:"required~Quantity is required"`
	Total_Price int `gorm:"not null" json:"total_price" valid:"required~Total Price is required"`
	UserID      uint
	ProductID   uint
	User        User
	Product     Product
}

func (p *TransactionHistory) Validate() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
