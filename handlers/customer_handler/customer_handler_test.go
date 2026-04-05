package handlers

import (
	"bytes"
	"context"
	customerrequest "e-commerce/dto/request/customer_request"
	customerresponse "e-commerce/dto/response/customer_response"
	"e-commerce/mocks/servicemocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CustomerHandler(t *testing.T) {
	testCases := []struct {
		description	string
		signupReq interface{}
		mockSetup func(m *servicemocks.CustomerSignupService)
		expectedStatus int
		expectedResponse customerresponse.CustomerSignupResponse
	}{
		{
			description: "Should return status OK and customer account created successfully message when all the details in the request (name , email , phone_number and password) are valid",
			signupReq: customerrequest.CustomerSignupRequest{
				Name : "rahul",
				Email : "rahul@gmail.com",
				PhoneNumber: "+917665718067",
				Password: "Rahul@1234",
			},
			mockSetup: func(m *servicemocks.CustomerSignupService) {
				m.On("CustomerSignup",mock.AnythingOfType("customerrequest.CustomerSignupRequest"),context.Background()).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: customerresponse.CustomerSignupResponse{
				IsSignUpSuccessful: true,
				Message: "Customer Account created successfully",
			}, 
		},
		{
			description: "Should return status Bad Request when email is not valid",
			signupReq: customerrequest.CustomerSignupRequest{
				Name : "rahul",
				Email : "rahulgmail.com",
				PhoneNumber: "+917665718067",
				Password: "Rahul@1234",
			},
			mockSetup: func(m *servicemocks.CustomerSignupService) {
				
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: customerresponse.CustomerSignupResponse{
				IsSignUpSuccessful: false,
			}, 
		},
		{
			description: "Should return status Bad Request when phone number is not valid",
			signupReq: customerrequest.CustomerSignupRequest{
				Name : "rahul",
				Email : "rahul@gmail.com",
				PhoneNumber: "+9176657180671234",
				Password: "Rahul@1234",
			},
			mockSetup: func(m *servicemocks.CustomerSignupService) {
				
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: customerresponse.CustomerSignupResponse{
				IsSignUpSuccessful: false,
			}, 
		},
		{
			description: "Should return status Bad Request when password is not valid",
			signupReq: customerrequest.CustomerSignupRequest{
				Name : "rahul",
				Email : "rahulgmail.com",
				PhoneNumber: "+917665718067",
				Password: "1234",
			},
			mockSetup: func(m *servicemocks.CustomerSignupService) {
				
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: customerresponse.CustomerSignupResponse{
				IsSignUpSuccessful: false,
			}, 
		},
		{
			description: "Should return status Internal Server Error and when all the details in the request (name , email , phone_number and password) are valid but there is an error in the service",
			signupReq: customerrequest.CustomerSignupRequest{
				Name : "rahul",
				Email : "rahul@gmail.com",
				PhoneNumber: "+917665718067",
				Password: "Rahul@1234",
			},
			mockSetup: func(m *servicemocks.CustomerSignupService) {
				m.On("CustomerSignup",mock.AnythingOfType("customerrequest.CustomerSignupRequest"),context.Background()).Return(errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: customerresponse.CustomerSignupResponse{
				IsSignUpSuccessful: false,
			}, 
		},
	}
	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			mockService := new(servicemocks.CustomerSignupService)
			tC.mockSetup(mockService)

			customerHandler := NewCustomerHandler(mockService)
			router := gin.New()
			router.POST("/customer/signup",customerHandler.CustomerSignup)

			var bodyBytes []byte
			switch v := tC.signupReq.(type){
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes , _ = json.Marshal(v)
			}

			req := httptest.NewRequest(http.MethodPost,"/customer/signup",bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type","application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w,req)
			assert.Equal(t,tC.expectedStatus,w.Code)

			var actualResponse customerresponse.CustomerSignupResponse
			err := json.NewDecoder(w.Body).Decode(&actualResponse)
			assert.NoError(t,err)
			assert.Equal(t,tC.expectedResponse.IsSignUpSuccessful,actualResponse.IsSignUpSuccessful)
			if tC.expectedResponse.Message != ""{
				assert.Equal(t,tC.expectedResponse.Message,actualResponse.Message)
			}

		})
	}
}