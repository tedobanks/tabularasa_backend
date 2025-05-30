package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/tedobanks/tabularasa_backend/db/sqlc"
)

// Server serves HTTP requests for our application.
type Server struct {
	store  *db.Queries // Use *db.Queries directly since your sqlc generated db.go provides it
	router *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(store *db.Queries) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Register your API routes here
	// Example:
	router.POST("/user", server.createUser)
	router.GET("/users", server.listUsers)
	router.GET("/user/:id", server.getUser)
	router.DELETE("/user/:id", server.deleteUser)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse is a helper for returning JSON errors
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
