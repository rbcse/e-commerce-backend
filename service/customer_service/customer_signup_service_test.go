package customerservice

import (
	"context"
	customerrequest "e-commerce/dto/request/customer_request"
	"e-commerce/mocks/repomocks"
	"e-commerce/mocks/servicemocks"
	"e-commerce/model"
	"testing"

	ae "e-commerce/error"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Test_CustomerSignupService(t *testing.T) {
	testCases := []struct {
		description string
		signupReq   customerrequest.CustomerSignupRequest
		expectedErr error
		mockSetup   func(m *repomocks.CustomerSignupRepository, h *servicemocks.PasswordHasher)
	}{
		{
			description: "Should return no error when all details(name , email , phone number and password are valid and email and phone number are unique) and an account for a customer should be created",
			signupReq: customerrequest.CustomerSignupRequest{
				Name:        "rahul_jain",
				Email:       "rahul@gmail.com",
				PhoneNumber: "+917665718067",
				Password:    "Rahul@1234",
			},
			mockSetup: func(m *repomocks.CustomerSignupRepository, h *servicemocks.PasswordHasher) {
				m.On("FindByEmail", context.Background(), "rahul@gmail.com").Return(nil, nil)
				m.On("FindByPhoneNumber", context.Background(), "+917665718067").Return(nil, nil)
				m.On("CreateCustomerAccount", mock.Anything).Return(nil)
				h.On("Hash", "Rahul@1234").Return("hash#1", nil)
			},
		},
		{
			description: "Should return Email Already exists error when all details(name , email , phone number and password are valid and email of customer already exists)",
			signupReq: customerrequest.CustomerSignupRequest{
				Name:        "rahul_jain",
				Email:       "rahul@gmail.com",
				PhoneNumber: "+917665718067",
				Password:    "Rahul@1234",
			},
			mockSetup: func(m *repomocks.CustomerSignupRepository, h *servicemocks.PasswordHasher) {
				m.On("FindByEmail", context.Background(), "rahul@gmail.com").Return(&model.Customer{}, nil)
			},
			expectedErr: ae.CustomerEmailAlreadyExists,
		},
		{
			description: "Should return Phone Number Already exists error when all details(name , email , phone number and password are valid and phone number of customer already exists)",
			signupReq: customerrequest.CustomerSignupRequest{
				Name:        "rahul_jain",
				Email:       "rahul@gmail.com",
				PhoneNumber: "+917665718067",
				Password:    "Rahul@1234",
			},
			mockSetup: func(m *repomocks.CustomerSignupRepository, h *servicemocks.PasswordHasher) {
				m.On("FindByEmail", context.Background(), "rahul@gmail.com").Return(nil, nil)
				m.On("FindByPhoneNumber", context.Background(), "+917665718067").Return(&model.Customer{}, nil)
			},
			expectedErr: ae.CustomerPhoneNumberAlreadyExists,
		},
		{
			description: "Should return Hashing error when all details(name , email , phone number and password are valid and email and phone number are unique but there is an error in hashing the password)",
			signupReq: customerrequest.CustomerSignupRequest{
				Name:        "rahul_jain",
				Email:       "rahul@gmail.com",
				PhoneNumber: "+917665718067",
				Password:    "Rahul@1234",
			},
			mockSetup: func(m *repomocks.CustomerSignupRepository, h *servicemocks.PasswordHasher) {
				m.On("FindByEmail", context.Background(), "rahul@gmail.com").Return(nil, nil)
				m.On("FindByPhoneNumber", context.Background(), "+917665718067").Return(nil, nil)
				h.On("Hash", "Rahul@1234").Return("", ae.HashingError)
			},
			expectedErr: ae.HashingError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			mockRepo := new(repomocks.CustomerSignupRepository)
			mockHasher := new(servicemocks.PasswordHasher)
			mockOTPService := new(servicemocks.OTPService)
			svc := NewCustomerSignupService(mockRepo, mockHasher , mockOTPService)
			tC.mockSetup(mockRepo, mockHasher)
			err := svc.CustomerSignup(tC.signupReq, context.Background())
			assert.Equal(t, err, tC.expectedErr)
		})
	}
}

func Test_PasswordHasher(t *testing.T) {
	hasher := &BcryptPasswordHasher{}
	password := "Rahul@1234"
	hashedPassword, err := hasher.Hash(password)
	assert.NoError(t, err)
	assert.NotEqual(t, hashedPassword, password)
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	assert.NoError(t, err)
}


func Test_VerifyCustomerOTP(t *testing.T) {
	testCases := []struct {
		description	string
		mockSetup func(os *servicemocks.OTPService , cr *repomocks.CustomerSignupRepository)
		identifier string
		otpType string
		otp string
		expectedErr error
	}{
		{
			description: "Should return no error when the otp entered by customer matches the correct otp and email should be marked as verified.",
			mockSetup: func(os *servicemocks.OTPService, cr *repomocks.CustomerSignupRepository) {
				os.On("VerifyOTP","rahul@gmail.com","EMAIL","123456").Return(nil)
				cr.On("MarkEmailVerified",mock.Anything,"rahul@gmail.com").Return(nil)
			},
			identifier: "rahul@gmail.com",
			otpType: "EMAIL",
			otp : "123456",
		},
		{
			description: "Should return Wrong OTP Error when the otp entered by customer does not matches the correct otp",
			mockSetup: func(os *servicemocks.OTPService, cr *repomocks.CustomerSignupRepository) {
				os.On("VerifyOTP","rahul@gmail.com","EMAIL","123456").Return(ae.ErrWrongOTPEntered)
			},
			identifier: "rahul@gmail.com",
			otpType: "EMAIL",
			otp : "123456",
			expectedErr: ae.ErrWrongOTPEntered,
		},
		{
			description: "Should return Record not found Error when the customer email does not exists in the DB",
			mockSetup: func(os *servicemocks.OTPService, cr *repomocks.CustomerSignupRepository) {
				os.On("VerifyOTP","rahul@gmail.com","EMAIL","123456").Return(nil)
				cr.On("MarkEmailVerified",mock.Anything,"rahul@gmail.com").Return(gorm.ErrRecordNotFound)
			},
			identifier: "rahul@gmail.com",
			otpType: "EMAIL",
			otp : "123456",
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			description: "Should return no error when the otp entered by customer matches the correct otp and phone Number should be marked as verified.",
			mockSetup: func(os *servicemocks.OTPService, cr *repomocks.CustomerSignupRepository) {
				os.On("VerifyOTP","+917665718067","PHONE","123456").Return(nil)
				cr.On("MarkPhoneNumberVerified",mock.Anything,"+917665718067").Return(nil)
			},
			identifier: "+917665718067",
			otpType: "PHONE",
			otp : "123456",
		},
		{
			description: "Should return Wrong otp error when the otp entered by customer does not matches the correct otp.",
			mockSetup: func(os *servicemocks.OTPService, cr *repomocks.CustomerSignupRepository) {
				os.On("VerifyOTP","+917665718067","PHONE","123456").Return(ae.ErrWrongOTPEntered)
			},
			identifier: "+917665718067",
			otpType: "PHONE",
			otp : "123456",
			expectedErr: ae.ErrWrongOTPEntered,
		},
		{
			description: "Should return record not found error when the customer's phone number does not exists in the DB",
			mockSetup: func(os *servicemocks.OTPService, cr *repomocks.CustomerSignupRepository) {
				os.On("VerifyOTP","+917665718067","PHONE","123456").Return(nil)
				cr.On("MarkPhoneNumberVerified",mock.Anything,"+917665718067").Return(gorm.ErrRecordNotFound)
			},
			identifier: "+917665718067",
			otpType: "PHONE",
			otp : "123456",
			expectedErr: gorm.ErrRecordNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			mockOTPService := new(servicemocks.OTPService)
			mockHasher := new(servicemocks.PasswordHasher)
			mockRepo := new(repomocks.CustomerSignupRepository)
			tC.mockSetup(mockOTPService,mockRepo)
			svc := NewCustomerSignupService(mockRepo,mockHasher,mockOTPService);
			err := svc.VerifyCustomerOTP(context.Background(),tC.identifier,tC.otpType,tC.otp);
			assert.Equal(t,tC.expectedErr,err);
		})
	}
}