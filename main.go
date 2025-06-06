package main

import (
	"database/sql"
	"log"

	"github.com/tedobanks/tabularasa_backend/api"
	db "github.com/tedobanks/tabularasa_backend/db/sqlc"
	"github.com/tedobanks/tabularasa_backend/routes"
	"github.com/tedobanks/tabularasa_backend/util"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Load configuration from .env or environment variables
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Open database connection
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close() // Ensure the connection is closed when main exits

	// Ping the database to verify the connection
	err = conn.Ping()
	if err != nil {
		log.Fatal("cannot ping db:", err)
	}
	log.Println("Successfully connected to the database!")

	// Create a new sqlc Queries instance
	store := db.New(conn)

	// Create a new Gin server and pass the store
	server := api.NewServer(store)

	// --- NEW: Register the API routes ---
	// Get the Gin router instance from the server
	router := server.Router

	// Register user routes
	routes.RegisterUserRoutes(router, server)
	// Register venue routes (if you have them implemented in routes/venue_routes.go)
	routes.RegisterVenueRoutes(router, server)
	// --- END NEW ---

	// Start the HTTP server
	log.Printf("Starting server at %s", config.ServerAddress)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
