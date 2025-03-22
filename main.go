package main

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"raya/config"
	"raya/database"
	"raya/routes"
)

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