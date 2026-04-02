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
)

func Test_CustomerSignupService(t *testing.T) {
	testCases := []struct {
		description	string
		signupReq customerrequest.CustomerSignupRequest
		expectedErr error
		mockSetup func(m *repomocks.CustomerSignupRepository , h *servicemocks.PasswordHasher)
	}{
		{
			description: "Should return no error when all details(name , email , phone number and password are valid and email and phone number are unique) and an account for a customer should be created",
			signupReq: customerrequest.CustomerSignupRequest{
				Name : "rahul_jain",
				Email: "rahul@gmail.com",
				PhoneNumber: "+917665718067",
				Password: "Rahul@1234",
			},
			mockSetup: func(m *repomocks.CustomerSignupRepository , h *servicemocks.PasswordHasher) {
				m.On("FindByEmail",context.Background(),"rahul@gmail.com").Return(nil,nil)
				m.On("FindByPhoneNumber",context.Background(),"+917665718067").Return(nil,nil)
				m.On("CreateCustomerAccount",mock.Anything).Return(nil)
				h.On("Hash","Rahul@1234").Return("hash#1",nil)
			},
		},
		{
			description: "Should return Email Already exists error when all details(name , email , phone number and password are valid and email of customer already exists)",
			signupReq: customerrequest.CustomerSignupRequest{
				Name : "rahul_jain",
				Email: "rahul@gmail.com",
				PhoneNumber: "+917665718067",
				Password: "Rahul@1234",
			},
			mockSetup: func(m *repomocks.CustomerSignupRepository , h *servicemocks.PasswordHasher) {
				m.On("FindByEmail",context.Background(),"rahul@gmail.com").Return(&model.Customer{},nil)
			},
			expectedErr: ae.CustomerEmailAlreadyExists,
		},
		{
			description: "Should return Phone Number Already exists error when all details(name , email , phone number and password are valid and phone number of customer already exists)",
			signupReq: customerrequest.CustomerSignupRequest{
				Name : "rahul_jain",
				Email: "rahul@gmail.com",
				PhoneNumber: "+917665718067",
				Password: "Rahul@1234",
			},
			mockSetup: func(m *repomocks.CustomerSignupRepository , h *servicemocks.PasswordHasher) {
				m.On("FindByEmail",context.Background(),"rahul@gmail.com").Return(nil,nil)
				m.On("FindByPhoneNumber",context.Background(),"+917665718067").Return(&model.Customer{},nil)
			},
			expectedErr: ae.CustomerPhoneNumberAlreadyExists,
		},
		{
			description: "Should return Hashing error when all details(name , email , phone number and password are valid and email and phone number are unique but there is an error in hashing the password)",
			signupReq: customerrequest.CustomerSignupRequest{
				Name : "rahul_jain",
				Email: "rahul@gmail.com",
				PhoneNumber: "+917665718067",
				Password: "Rahul@1234",
			},
			mockSetup: func(m *repomocks.CustomerSignupRepository , h *servicemocks.PasswordHasher) {
				m.On("FindByEmail",context.Background(),"rahul@gmail.com").Return(nil,nil)
				m.On("FindByPhoneNumber",context.Background(),"+917665718067").Return(nil,nil)
				h.On("Hash","Rahul@1234").Return("",ae.HashingError)
			},
			expectedErr: ae.HashingError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			mockRepo := new(repomocks.CustomerSignupRepository)
			mockHasher := new(servicemocks.PasswordHasher)
			svc := NewCustomerSignupService(mockRepo,mockHasher)
			tC.mockSetup(mockRepo,mockHasher)
			err := svc.CustomerSignup(tC.signupReq,context.Background())
			assert.Equal(t,err,tC.expectedErr)
		})
	}
}

func Test_PasswordHasher(t *testing.T) {
	hasher := &BcryptPasswordHasher{}
	password := "Rahul@1234"
	hashedPassword , err := hasher.Hash(password)
	assert.NoError(t,err)
	assert.NotEqual(t,hashedPassword,password)
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(password))
	assert.NoError(t,err)
}	