package customerservice

import (
	"context"
	customerrequest "e-commerce/dto/request/customer_request"
	ae "e-commerce/error"
	"e-commerce/model"
	customerrepository "e-commerce/repository/customer_repository"

	"golang.org/x/crypto/bcrypt"
)

type CustomerSignupService interface {
	CustomerSignup(req customerrequest.CustomerSignupRequest , ctx context.Context) error
}

type customerSignupService struct {
	repo customerrepository.CustomerSignupRepository
}

func NewCustomerSignupService(repo customerrepository.CustomerSignupRepository) CustomerSignupService{
	return &customerSignupService{
		repo : repo,
	}
}

func (cs *customerSignupService) CustomerSignup(req customerrequest.CustomerSignupRequest , ctx context.Context) error {

	existingCustomerByEmail , _ := cs.repo.FindByEmail(ctx,req.Email)
	if existingCustomerByEmail != nil {
		return ae.CustomerEmailAlreadyExists
	}

	existingCustomerByPhoneNumber , _ := cs.repo.FindByPhoneNumber(ctx,req.PhoneNumber)
	if existingCustomerByPhoneNumber != nil {
		return ae.CustomerPhoneNumberAlreadyExists
	}

	hashedPassword , err := bcrypt.GenerateFromPassword([]byte(req.Password),bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	newCustomer := model.NewCustomer(req.Name , req.Email , req.PhoneNumber , string(hashedPassword))

	err = cs.repo.CreateCustomerAccount(newCustomer)

	return err

}