package routes

import (
	"raya/controllers"
	"raya/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

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
		api.GET("/links", controllers.GetAllCategories)

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
			
			// Link management dalam kategori
			admin.GET("/categories/:category_id/links", controllers.GetLinksByCategory)
			admin.POST("/categories/:category_id/links", controllers.CreateLink)
			admin.PATCH("/categories/:category_id/links/:link_id", controllers.UpdateLink)
			admin.DELETE("/categories/:category_id/links/:link_id", controllers.DeleteLink)

			// Category management
			admin.GET("/categories", controllers.GetAllCategories)
			admin.GET("/category/:id", controllers.GetCategoryByID)
			admin.POST("/category", controllers.CreateCategory)
			admin.PATCH("/category/:id", controllers.UpdateCategory)
			admin.DELETE("/category/:id", controllers.DeleteCategory)
		}
	}
	return r
}