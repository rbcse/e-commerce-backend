package server

import (
    "e-commerce/config"
    "e-commerce/db"
    "e-commerce/routes"
    "github.com/gin-gonic/gin"
)

func Start() {
    cfg := config.Load()
    database := db.Connect(cfg.DatabaseURL)

    r := gin.Default()
    routes.Register(r, database)
    r.Run(":" + cfg.Port)
}