package entity

import "github.com/asaskevich/govalidator"

type Category struct {
	GormModel
	Type                string    `gorm:"not null" json:"type" valid:"required~Type is required"`
	Sold_Product_Amount int       `gorm:"not null" json:"sold_product_amount"`
	Product             []Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"product"`
}

func (p *Category) Validate() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
