package server

import (
	"e-commerce/config"
	"e-commerce/db"
	"e-commerce/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start() {
	cfg := config.Load()
	database := db.Connect(cfg.DatabaseURL)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	routes.Register(r, database)
	r.Run(":" + cfg.Port)
}
