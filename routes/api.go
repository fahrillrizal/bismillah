package routes

import (
	"raya/controllers"
	"raya/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{
			"https://sekawan-grup.com",
			"https://api.sekawan-grup.com",
			"http://localhost:3000",
			"http://localhost:8080",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{
			"Origin", 
			"Content-Type", 
			"Content-Length", 
			"Accept-Encoding", 
			"X-CSRF-Token", 
			"Authorization",
			"Accept",
		},
		ExposeHeaders:    []string{
			"Content-Length", 
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	r.Use(func(c *gin.Context) {
		if db == nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Database connection error"})
			return
		}
		c.Set("db", db)
		c.Next()
	})

	api := r.Group("/api")
	{
		api.GET("/categories-with-links", controllers.GetCategoriesWithLinks)//untuk section service

		// Auth
		api.POST("/login", controllers.LoginUser)

		admin := api.Group("/")
		admin.Use(middleware.AuthMiddleware())
		{
			admin.PATCH("/change-password", controllers.ChangePassword)
			admin.POST("/logout", controllers.LogoutUser)

			// Link management
			admin.GET("/links/all", controllers.GetAllLinks)
			admin.GET("/links/:id", controllers.GetLinkByID)
			admin.GET("/links", controllers.GetLinks)//untuk dashboard/links
			
			// Link management dalam kategori
			admin.GET("/categories/:category_id/links", controllers.GetLinksByCategory)
			admin.POST("/categories/:category_id/links", controllers.CreateLink)
			admin.PATCH("/categories/:category_id/links/:link_id", controllers.UpdateLink)
			admin.DELETE("/categories/:category_id/links/:link_id", controllers.DeleteLink)

			// Category management
			admin.GET("/categories", controllers.GetCategories)//untuk dashboard/categories
			admin.GET("/category/:id", controllers.GetCategoryByID)
			admin.POST("/category", controllers.CreateCategory)
			admin.PATCH("/category/:id", controllers.UpdateCategory)
			admin.DELETE("/category/:id", controllers.DeleteCategory)
		}
	}
	return r
}