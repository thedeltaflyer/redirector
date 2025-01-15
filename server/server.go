package server

import (
	"github.com/gin-gonic/gin"

	"redirector/controllers"
	"redirector/database"
	"redirector/middleware"
	"redirector/models"
)

// Run starts the HTTP server with the specified binding address and debug mode settings.
// It initializes controllers, middleware, and routes for handling HTTP requests.
// The function panics if the server fails to start.
func Run(bind string, debug bool) {
	// Set ReleaseMode if we're not debugging.
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create a basic Gin Engine.
	r := gin.New()

	// Add the Logging and Recovery middleware.
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// KV for the "redirects" bucket.
	redirectKV := &models.KVWrapper{
		DB:     database.GetDB(),
		Bucket: []byte("redirects"),
	}

	// KV for the "api_keys" bucket.
	apiKeyKV := &models.KVWrapper{
		DB:     database.GetDB(),
		Bucket: []byte("api_keys"),
	}

	// KV for the "health_checks" bucket.
	healthKV := &models.KVWrapper{
		DB:     database.GetDB(),
		Bucket: []byte("health_checks"),
	}

	// Create the controllers.
	root := &controllers.RootController{}
	health := &controllers.HealthController{
		KV: healthKV,
	}
	redirector := &controllers.RedirectorController{
		KV: redirectKV,
	}

	// Set up static and health routes
	rootGroup := r.Group("/")
	rootGroup.GET("", root.HandleGet)
	rootGroup.GET("/health", health.HandleGet)

	// Set up unauthenticated redirection routes
	redirectorGroup := r.Group("/")
	redirectorGroup.GET("/:key/*mode", redirector.HandleGet)

	// Set up authenticated redirection routes
	createRedirectorGroup := r.Group("/")
	createRedirectorGroup.Use(middleware.TokenAuthMiddleware(apiKeyKV))
	createRedirectorGroup.POST("", redirector.HandlePost)
	createRedirectorGroup.POST("/:key", redirector.HandlePost)
	createRedirectorGroup.PUT("/:key", redirector.HandlePutWithKey)

	// Start the server
	err := r.Run(bind)
	if err != nil {
		panic(err)
	}
}
