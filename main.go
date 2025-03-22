package main

import (
	"log"
	"os"
	_ "raya/docs"
	"github.com/joho/godotenv"
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


	router.Run(":" + os.Getenv("PORT"))
}