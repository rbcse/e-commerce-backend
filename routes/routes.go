package routes

import (
	"e-commerce/handlers"
	customerrepository "e-commerce/repository/customer_repository"
	customerservice "e-commerce/service/customer_service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api/v1")
	customerRepo := customerrepository.NewCustomerSignupRepository(db)
	customerService := customerservice.NewCustomerSignupService(customerRepo)
	handlers.RegisterCustomerRoutes(api, customerService)
}
