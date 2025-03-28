package main

import (
	"log"
	"os"
	"time"

	_ "raya/docs"
	
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"raya/config"
	"raya/database"
	"raya/routes"
)

// @title Raya API
// @version 1.0
// @description API untuk aplikasi Raya menggunakan Gin framework
// @host api.sekawan-grup.com
// @BasePath /api
func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	
	if err := database.SeedDatabase(db); err != nil {
		log.Printf("Error seeding database: %v", err)
	}

	router := routes.SetupRouter(db)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{
			"https://sekawan-grup.com",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{
			"Origin", 
			"Content-Type", 
			"Authorization", 
			"Accept",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		c.Status(200)
	})

	router.Run(":" + os.Getenv("PORT"))
}