package server

import (
	"context"
	"e-commerce/config"
	"e-commerce/db"
	"e-commerce/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"e-commerce/common/logger"
)

func Start() {
	logger.InitLogger(logger.Config{
		Level: "info",
	})
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		logger.Fatal("failed to connect to redis: %v", err)
	}
	cfg := config.Load()
	database := db.Connect(cfg.DatabaseURL)

	logger.Info("database connected successfully");

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	routes.Register(r, database, redisClient)
	
	logger.Info("server starting on port %s",cfg.Port);
	r.Run(":" + cfg.Port)

}
