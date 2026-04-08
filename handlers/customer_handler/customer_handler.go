package handlers

import (
	"e-commerce/app"
	customerrequest "e-commerce/dto/request/customer_request"
	customerresponse "e-commerce/dto/response/customer_response"
	customerservice "e-commerce/service/customer_service"
	otpservice "e-commerce/service/otp_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	customerService customerservice.CustomerSignupService
	otpService      otpservice.OTPService
}

func NewCustomerHandler(
	customerSvc customerservice.CustomerSignupService,
	otpSvc otpservice.OTPService,
) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerSvc,
		otpService:      otpSvc,
	}
}

func RegisterCustomerRoutes(
	rg *gin.RouterGroup,
	customerSvc customerservice.CustomerSignupService,
	otpSvc otpservice.OTPService,
) {
	h := NewCustomerHandler(customerSvc, otpSvc)
	customers := rg.Group("/customer")
	{
		customers.POST(app.CustomerSignupEndPoint, h.CustomerSignup)
		customers.POST(app.GenerateOTP, h.GenerateOTP)
	}
}

func (ch *CustomerHandler) CustomerSignup(c *gin.Context) {
	var req customerrequest.CustomerSignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, customerresponse.CustomerSignupResponse{
			IsSignUpSuccessful: false,
			Message:            err.Error(),
		})
		return
	}

	if err := ch.customerService.CustomerSignup(req, c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, customerresponse.CustomerSignupResponse{
			IsSignUpSuccessful: false,
			Message:            err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, customerresponse.CustomerSignupResponse{
		IsSignUpSuccessful: true,
		Message:            "Customer Account created successfully",
	})
}

func (ch *CustomerHandler) GenerateOTP(c *gin.Context) {
	var req customerrequest.GenerateOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, customerresponse.GenerateOTPResponse{
			Message: err.Error(),
		})
		return
	}

	otp, err := ch.otpService.GenerateOTP(req.Identifier, req.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, customerresponse.GenerateOTPResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, customerresponse.GenerateOTPResponse{
		Otp:     otp,
		Message: "OTP generated successfully",
	})
}
