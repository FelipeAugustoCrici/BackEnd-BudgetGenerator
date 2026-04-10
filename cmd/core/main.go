package main

import (
	"log"
	"os"

	"budgetgen/internal/auth"
	"budgetgen/internal/db"
	"budgetgen/internal/handler/corehandler"
	"budgetgen/internal/handler/crmhandler"
	"budgetgen/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, using environment variables")
	}

	db.Connect()
	storage.Connect()

	r := gin.Default()
	r.Use(corsMiddleware())
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	api := r.Group("/api", auth.Middleware())
	{
		api.GET("/me", corehandler.Me)

		api.GET("/quotes", corehandler.ListQuotes)
		api.GET("/quotes/:id", corehandler.GetQuote)
		api.POST("/quotes", corehandler.CreateQuote)
		api.PUT("/quotes/:id", corehandler.UpdateQuote)
		api.DELETE("/quotes/:id", corehandler.DeleteQuote)

		api.GET("/templates", corehandler.ListTemplates)
		api.GET("/templates/:id", corehandler.GetTemplate)
		api.POST("/templates", corehandler.CreateTemplate)
		api.PUT("/templates/:id", corehandler.UpdateTemplate)
		api.DELETE("/templates/:id", corehandler.DeleteTemplate)

		api.GET("/settings", corehandler.GetSettings)
		api.PUT("/settings", corehandler.UpsertSettings)

		api.POST("/upload", corehandler.Upload)
		api.POST("/upload/presign", corehandler.PresignUpload)
		api.GET("/image-proxy", corehandler.ImageProxy)
		api.POST("/ai/quote", corehandler.GenerateQuote)

		// CRM - Clients
		api.GET("/clients", crmhandler.ListClients)
		api.POST("/clients", crmhandler.CreateClient)
		api.GET("/clients/:id", crmhandler.GetClient)
		api.PUT("/clients/:id", crmhandler.UpdateClient)
		api.DELETE("/clients/:id", crmhandler.DeleteClient)

		// CRM - Contracts
		api.GET("/contracts", crmhandler.ListContracts)
		api.POST("/contracts", crmhandler.CreateContract)
		api.GET("/contracts/:id", crmhandler.GetContract)
		api.PUT("/contracts/:id", crmhandler.UpdateContract)
		api.POST("/contracts/:id/send", crmhandler.SendContract)
		api.POST("/contracts/:id/view", crmhandler.ViewContract)
		api.POST("/contracts/:id/sign", crmhandler.SignContract)
		api.POST("/contracts/:id/refuse", crmhandler.RefuseContract)
		api.GET("/contracts/:id/events", crmhandler.ListContractEvents)
		api.GET("/contracts/by-budget/:budgetId", crmhandler.GetContractByBudget)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("CORE_PORT")
	}
	if port == "" {
		port = "9000"
	}
	log.Printf("core-service running on :%s", port)
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
