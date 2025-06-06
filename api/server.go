package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/tedobanks/tabularasa_backend/db/sqlc"
	// "github.com/tedobanks/tabularasa_backend/routes"
)

// Server serves HTTP requests for our application.
type Server struct {
	store  *db.Queries
	Router *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(store *db.Queries) *Server {
	server := &Server{store: store}
	server.Router = gin.Default()

	// Register all your API routes here by calling functions from the routes package.
	// This keeps your server.go clean and delegates route definition to dedicated files.
	// routes.RegisterUserRoutes(router, server)
	// routes.RegisterVenueRoutes(router, server)

	// If you had other resources (e.g., products), you would add calls here:
	// routes.RegisterProductRoutes(router, server)
	// routes.RegisterOrderRoutes(router, server)

	// server.router = router
	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}

// errorResponse is a helper for returning JSON errors
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// --- Placeholder methods for your handlers ---
// In a real application, these methods would contain your actual business logic,
// interacting with server.store. They would typically be in separate handler files
// (e.g., api/user_handlers.go) but are included here for a complete example.

// CreateUser handles the creation of a new user.
