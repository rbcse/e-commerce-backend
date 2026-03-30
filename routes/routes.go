package routes

import (
    "e-commerce/handlers"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func Register(r *gin.Engine, db *gorm.DB) {
    api := r.Group("/api/v1")
	handlers.RegisterCustomerRoutes(api,db)
}