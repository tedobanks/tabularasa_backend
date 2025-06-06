package routes // You might put this in a 'routes' package or similar

import (
	"github.com/gin-gonic/gin"
	"github.com/tedobanks/tabularasa_backend/api" // Import your api package
)

// RegisterUserRoutes registers all user-related API routes.
// It takes a *gin.Engine (the main router) and an *api.Server instance.
func RegisterUserRoutes(router *gin.Engine, server *api.Server) {

	// Public routes (no authentication required)
	router.POST("/register", server.CreateUser) // User registration
	router.POST("/login", server.LoginUser)     // User login

	// Authenticated routes (require a valid JWT token)
	// All routes within this group will be protected by the AuthMiddleware
	// Note: We use 'router' here, not 'server.Router', because this function
	// is passed the main Gin engine.
	userRoutes := router.Group("/user").Use(api.AuthMiddleware()) // Use api.AuthMiddleware()
	{
		userRoutes.POST("/logout", server.LogoutUser) // New: Logout endpoint
		userRoutes.GET("/:id", server.GetUser)        // Path becomes /user/:id
		userRoutes.GET("/", server.ListUsers)         // Path becomes /user
		userRoutes.PUT("/:id", server.UpdateUser)
		userRoutes.DELETE("/:id", server.DeleteUser)
	}
}
