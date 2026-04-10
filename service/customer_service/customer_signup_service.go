package customerservice

import (
	"context"
	customerrequest "e-commerce/dto/request/customer_request"
	ae "e-commerce/error"
	"e-commerce/model"
	customerrepository "e-commerce/repository/customer_repository"
	otpservice "e-commerce/service/otp_service"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type BcryptPasswordHasher struct{}

func (b *BcryptPasswordHasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

type CustomerSignupService interface {
	CustomerSignup(req customerrequest.CustomerSignupRequest, ctx context.Context) error
	VerifyCustomerOTP(ctx context.Context , identifier , otpType , otp string) error
}

type customerSignupService struct {
	repo   customerrepository.CustomerSignupRepository
	hasher PasswordHasher
	otpService otpservice.OTPService
}

func NewCustomerSignupService(repo customerrepository.CustomerSignupRepository, hasher PasswordHasher , otpService otpservice.OTPService) CustomerSignupService {
	return &customerSignupService{
		repo:   repo,
		hasher: hasher,
		otpService: otpService,
	}
}

func (cs *customerSignupService) CustomerSignup(req customerrequest.CustomerSignupRequest, ctx context.Context) error {

	existingCustomerByEmail, _ := cs.repo.FindByEmail(ctx, req.Email)
	if existingCustomerByEmail != nil {
		return ae.CustomerEmailAlreadyExists
	}

	existingCustomerByPhoneNumber, _ := cs.repo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if existingCustomerByPhoneNumber != nil {
		return ae.CustomerPhoneNumberAlreadyExists
	}

	hashedPassword, err := cs.hasher.Hash(req.Password)

	if err != nil {
		return err
	}

	newCustomer := model.NewCustomer(req.Name, req.Email, req.PhoneNumber, string(hashedPassword))

	err = cs.repo.CreateCustomerAccount(newCustomer)

	return err

}

func (cs *customerSignupService) VerifyCustomerOTP(ctx context.Context , identifier , otp_type , otp string) error {

	err := cs.otpService.VerifyOTP(identifier,otp_type,otp);
	if err != nil {
		return err
	}

	if otp_type == "EMAIL" {
		err = cs.repo.MarkEmailVerified(ctx,identifier)
	} else {
		err = cs.repo.MarkPhoneNumberVerified(ctx,identifier)
	}

	return err
}