package model

type Customer struct {
	CustomerId            uint   `gorm:"primaryKey;autoIncrement" json:"customer_id"`
	Name                  string `gorm:"type:VARCHAR(50);not null" json:"name"`
	Email                 string `gorm:"type:VARCHAR(255);UNIQUE;not null" json:"email"`
	PhoneNumber           string `gorm:"type:VARCHAR(14);UNIQUE;not null" json:"phone_number"`
	PasswordHash          string `gorm:"type:VARCHAR(255);not null" json:"password_hash"`
	IsEmailVerified       bool   `gorm:"type:BOOLEAN;default:false" json:"is_email_verfied"`
	IsPhoneNumberVerified bool   `gorm:"type:BOOLEAN;default:false" json:"is_phone_number_verfied"`
}

func NewCustomer(name, email, phone_number, hashed_password string) *Customer {
	return &Customer{
		Name:         name,
		Email:        email,
		PhoneNumber:  phone_number,
		PasswordHash: hashed_password,
	}
}
