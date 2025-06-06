package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tedobanks/tabularasa_backend/api"
)

func RegisterVenueRoutes(router *gin.Engine, server *api.Server) {
	venueRoutes := router.Group("/venues").Use(api.AuthMiddleware())
	{
		// POST /venues - Create a new user
		venueRoutes.POST("", server.CreateVenue)

		// GET /venues - List all users
		venueRoutes.GET("", server.ListVenues)

		// GET /venues/:id - Get a specific user by ID
		// Note: The ":id" is a path parameter. Gin automatically parses it.
		venueRoutes.GET("/:name", server.GetVenue)

		// DELETE /venues/:id - Delete a specific user by ID
		venueRoutes.DELETE("/:id", server.DeleteVenue)

		// Add more venue-specific routes here as needed, e.g.:
		venueRoutes.PUT("/:id", server.UpdateVenue) // For updating a venue
	}

	// If you prefer to have /venue for single resource operations (e.g., POST /venue),
	// you can register it directly on the main router, outside the /venues group:
	// router.POST("/venue", server.CreateVenue) // If you want this specific path
}
