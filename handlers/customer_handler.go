package handlers

import (
	"e-commerce/app"
	customerrequest "e-commerce/dto/request/customer_request"
	customerresponse "e-commerce/dto/response/customer_response"
	customerservice "e-commerce/service/customer_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service customerservice.CustomerSignupService
}

func NewCustomerHandler(svc customerservice.CustomerSignupService)*CustomerHandler{
	return  &CustomerHandler{
		service: svc,
	}
}

func RegisterCustomerRoutes(rg *gin.RouterGroup , svc customerservice.CustomerSignupService) {
	h := NewCustomerHandler(svc)
	customers := rg.Group("/customer")
	{
		customers.POST(app.CustomerSignupEndPoint, h.CustomerSignup)
	}
}

func (ch *CustomerHandler) CustomerSignup(c *gin.Context) {

	var req customerrequest.CustomerSignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,customerresponse.CustomerSignupResponse{
			IsSignUpSuccessful: false,
			Message: err.Error(),
		})
		return
	}

	err := ch.service.CustomerSignup(req,c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError , customerresponse.CustomerSignupResponse{
			IsSignUpSuccessful: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK , customerresponse.CustomerSignupResponse{
		IsSignUpSuccessful: true,
		Message: "Customer Account created successfully",
	}) 

}
