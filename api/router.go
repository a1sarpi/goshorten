package api

import (
	"github.com/a1sarpi/goshorten/api/handlers"
	"github.com/a1sarpi/goshorten/api/storage"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title GoShorten API
// @version 1.0
// @description A URL shortening service with support for PostgreSQL and in-memory storage
// @host localhost:8080
// @BasePath /
// @schemes http https
// @contact.name API Support
// @contact.url https://github.com/a1sarpi/goshorten
// @license.name MIT
// @license.url https://github.com/a1sarpi/goshorten/blob/main/LICENSE
func SetupRouter(store storage.Storage) *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	postHandler := handlers.NewPostHandler(store)
	getHandler := handlers.NewGetHandler(store)

	router.POST("/shorten", postHandler.HandleShorten)
	router.GET("/:shortcode", getHandler.HandleRedirect)

	return router
}
