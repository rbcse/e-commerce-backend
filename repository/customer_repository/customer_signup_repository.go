package customerrepository

import (
	"context"
	"e-commerce/model"

	"gorm.io/gorm"
)

type CustomerSignupRepository interface {
	CreateCustomerAccount(customer *model.Customer) error
	FindByEmail(ctx context.Context , email string) (*model.Customer , error)
	FindByPhoneNumber(ctx context.Context , phone_number string) (*model.Customer , error)
}

type customerSignupRepository struct {
	db *gorm.DB
}

func NewCustomerSignupRepository(db *gorm.DB) CustomerSignupRepository{
	return &customerSignupRepository{
		db : db,
	}
}

func (cr *customerSignupRepository) CreateCustomerAccount(customer *model.Customer) error {
	result := cr.db.Create(customer)
	return result.Error
}

func (cr *customerSignupRepository) FindByEmail(ctx context.Context , email string) (*model.Customer , error) {
	var customer *model.Customer
	result := cr.db.WithContext(ctx).Where("email = ?",email).First(&customer)
	if result.Error != nil{
		return nil , result.Error
	}
	return customer , nil
}

func (cr *customerSignupRepository) FindByPhoneNumber(ctx context.Context , phone_number string) (*model.Customer , error) {
	var customer *model.Customer
	result := cr.db.WithContext(ctx).Where("phone_number = ?",phone_number).First(&customer)
	if result.Error != nil {
		return nil , result.Error
	}
	return customer,nil
}