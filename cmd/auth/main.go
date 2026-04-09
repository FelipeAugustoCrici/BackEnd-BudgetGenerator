package main

import (
	"log"
	"os"

	"budgetgen/internal/db"
	"budgetgen/internal/handler/authhandler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, using environment variables")
	}

	db.Connect()

	r := gin.Default()
	r.Use(corsMiddleware())

	r.POST("/auth/register", authhandler.Register)
	r.POST("/auth/login", authhandler.Login)

	port := os.Getenv("AUTH_PORT")
	if port == "" {
		port = "9001"
	}
	log.Printf("auth-service running on :%s", port)
	r.Run(":" + port)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
