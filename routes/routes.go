package routes

import (
	customer_handler "e-commerce/handlers/customer_handler"
	customerrepository "e-commerce/repository/customer_repository"
	otprepository "e-commerce/repository/otp_repository"
	customerservice "e-commerce/service/customer_service"
	otpservice "e-commerce/service/otp_service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Register(r *gin.Engine, db *gorm.DB, redisClient *redis.Client) {
	api := r.Group("/api/v1")

	// Repositories
	customerRepo := customerrepository.NewCustomerSignupRepository(db)
	otpRepo := otprepository.NewOTPRepository(redisClient, 5*time.Minute)

	// Services
	hasher := &customerservice.BcryptPasswordHasher{}
	customerService := customerservice.NewCustomerSignupService(customerRepo, hasher)
	otpService := otpservice.NewOTPService(otpRepo)

	// Handlers
	customer_handler.RegisterCustomerRoutes(api, customerService, otpService)
}
