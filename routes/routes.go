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

	// Services Dependencies
	hasher := &customerservice.BcryptPasswordHasher{}
	otpGenerator := &otpservice.DefaultOTPGenerator{}
	senderFactory := &otpservice.DefaultSenderFactory{}

	// Services
	otpService := otpservice.NewOTPService(otpRepo,otpGenerator,senderFactory)
	customerService := customerservice.NewCustomerSignupService(customerRepo, hasher , otpService)

	// Handlers
	customer_handler.RegisterCustomerRoutes(api, customerService, otpService)
}
