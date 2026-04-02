package customerservice

import (
	"context"
	customerrequest "e-commerce/dto/request/customer_request"
	ae "e-commerce/error"
	"e-commerce/model"
	customerrepository "e-commerce/repository/customer_repository"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface{
	Hash(password string) (string,error)
}

type BcryptPasswordHasher struct{}

func (b *BcryptPasswordHasher) Hash(password string) (string,error){
	hashedPassword , err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	return string(hashedPassword) , err
}

type CustomerSignupService interface {
	CustomerSignup(req customerrequest.CustomerSignupRequest , ctx context.Context) error
}

type customerSignupService struct {
	repo customerrepository.CustomerSignupRepository
	hasher PasswordHasher
}

func NewCustomerSignupService(repo customerrepository.CustomerSignupRepository , hasher PasswordHasher) CustomerSignupService{
	return &customerSignupService{
		repo : repo,
		hasher: hasher,
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

	hashedPassword , err := cs.hasher.Hash(req.Password)

	if err != nil {
		return err
	}

	newCustomer := model.NewCustomer(req.Name , req.Email , req.PhoneNumber , string(hashedPassword))

	err = cs.repo.CreateCustomerAccount(newCustomer)

	return err

}