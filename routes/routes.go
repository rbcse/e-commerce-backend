package routes

import (
	customer_handler "e-commerce/handlers/customer_handler"
	customerrepository "e-commerce/repository/customer_repository"
	customerservice "e-commerce/service/customer_service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api/v1")
	customerRepo := customerrepository.NewCustomerSignupRepository(db)
	hasher := &customerservice.BcryptPasswordHasher{}
	customerService := customerservice.NewCustomerSignupService(customerRepo,hasher)
	customer_handler.RegisterCustomerRoutes(api, customerService)
}
