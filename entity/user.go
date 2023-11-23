package entity

import (
	"errors"
	"finalProject4/helper"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	GormModel
	Full_Name          string               `gorm:"not null" json:"full_name" valid:"required~Your username is required"`
	Email              string               `gorm:"not null" json:"email" valid:"required~Your email is required,email~Invalid email format"`
	Password           string               `gorm:"not null" json:"password" valid:"required~Your password is required,minstringlength(6)~Password has to have minimum length of 6 characters"`
	Role               string               `gorm:"not null" json:"role"`
	Balance            int                  `gorm:"not null;default:0" json:"balance" valid:"range(0|100000000)~Balance must be under 100000000"`
	TransactionHistory []TransactionHistory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"transcationhistory"`
}

func (u *User) Validate() error {
	if u.Role != "admin" && u.Role != "customer" {
		return errors.New("Role must be admin or customer")
	}
	_, err := govalidator.ValidateStruct(u)
	return err
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	u.Password = helper.HashPass(u.Password)
	err = nil
	return
}
