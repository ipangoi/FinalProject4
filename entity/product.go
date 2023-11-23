package entity

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

type Product struct {
	GormModel
	Title              string `gorm:"not null" json:"title" valid:"required~Title is required"`
	Price              int    `gorm:"not null" json:"price" valid:"required~Price is required,range(0|50000000)~Price must be under 50000000"`
	Stock              int    `gorm:"not null" json:"stock" valid:"required~Stock is required"`
	CategoryID         uint
	TransactionHistory []TransactionHistory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"transcationhistory"`
}

func (p *Product) Validate() error {
	if p.Stock < 5 {
		return errors.New("Stock must be at least 5")
	}

	_, err := govalidator.ValidateStruct(p)
	return err
}
