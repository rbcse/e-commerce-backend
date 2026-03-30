package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CustomerHandler struct {
	db *gorm.DB
}

func RegisterCustomerRoutes(rg *gin.RouterGroup , db *gorm.DB){
	h := &CustomerHandler{
		db : db,
	}

	customers := rg.Group("/customer")
	{
		customers.POST("/signup",h.CustomerSignup)
	}
}

func (h *CustomerHandler) CustomerSignup(c *gin.Context){}