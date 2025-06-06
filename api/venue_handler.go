package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/tedobanks/tabularasa_backend/db/sqlc" // Assuming this path is correct
)

// Helper function to create sql.NullString from a string
// func newNullString(s string) sql.NullString {
// 	if len(s) == 0 {
// 		return sql.NullString{Valid: false}
// 	}
// 	return sql.NullString{String: s, Valid: true}
// }

// Helper function to create sql.NullInt32 from an int
func newNullInt32(i int) sql.NullInt32 {
	if i == 0 {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: int32(i), Valid: true}
}

// Helper function to create sql.NullTime from a time.Time
func newNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}

// Helper function to create uuid.NullUUID from a uuid.UUID
func newNullUUID(u uuid.UUID) uuid.NullUUID {
	if u == uuid.Nil {
		return uuid.NullUUID{Valid: false}
	}
	return uuid.NullUUID{UUID: u, Valid: true}
}

// createVenueRequest defines the request body for creating a new venue.
type createVenueRequest struct {
	ImageLinks      []string  `json:"image_links"`
	Name            string    `json:"name" binding:"required"`
	Type            string    `json:"type"`
	Description     string    `json:"description"`
	Location        string    `json:"location" binding:"required"`
	Dimension       string    `json:"dimension"`
	Capacity        int       `json:"capacity"`
	Facilities      []string  `json:"facilities"`
	HasAccomodation bool      `json:"has_accomodation"`
	RoomType        string    `json:"room_type"`
	NoOfRooms       int       `json:"no_of_rooms"`
	Sleeps          string    `json:"sleeps"`
	BedType         string    `json:"bed_type"`
	Rent            int       `json:"rent"`
	OwnedBy         uuid.UUID `json:"owned_by" binding:"required"` // Assuming owned_by is provided by client
	IsAvailable     bool      `json:"is_available"`
	OpensAt         time.Time `json:"opens_at"`
	ClosesAt        time.Time `json:"closes_at"`
	RentalDays      string    `json:"rental_days"`
	BookingPrice    int       `json:"booking_price"`
}

// createVenue handles the creation of a new venue.
// POST /venues
func (server *Server) CreateVenue(ctx *gin.Context) {
	var req createVenueRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateVenueParams{
		Column1:         req.ImageLinks,
		Name:            req.Name,
		Type:            newNullString(req.Type),
		Description:     newNullString(req.Description),
		Location:        req.Location,
		Dimension:       newNullString(req.Dimension),
		Capacity:        newNullInt32(req.Capacity),
		Column8:         req.Facilities,
		HasAccomodation: sql.NullBool{Bool: req.HasAccomodation, Valid: true},
		RoomType:        newNullString(req.RoomType),
		NoOfRooms:       newNullInt32(req.NoOfRooms),
		Sleeps:          newNullString(req.Sleeps),
		BedType:         newNullString(req.BedType),
		Rent:            newNullInt32(req.Rent),
		OwnedBy:         newNullUUID(req.OwnedBy),
		IsAvailable:     sql.NullBool{Bool: req.IsAvailable, Valid: true},
		OpensAt:         newNullTime(req.OpensAt),
		ClosesAt:        newNullTime(req.ClosesAt),
		RentalDays:      newNullString(req.RentalDays),
		BookingPrice:    newNullInt32(req.BookingPrice),
	}

	venue, err := server.store.CreateVenue(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, venue)
}

// getVenueRequest defines the URI parameter for getting a venue by name.
type getVenueRequest struct {
	Name string `uri:"name" binding:"required"` // Changed to 'name' and removed 'uuid' binding
}

// getVenue handles fetching a single venue by Name.
// GET /venues/name/:name
func (server *Server) GetVenue(ctx *gin.Context) {
	var req getVenueRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Assuming you have a method in your store to get a venue by name
	venue, err := server.store.GetVenueByName(ctx, req.Name) // Changed to GetVenueByName
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("venue not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, venue)
}

// listVenuesRequest defines optional query parameters for listing venues (e.g., pagination).
type listVenuesRequest struct {
	// You could add Limit int32 `form:"limit" binding:"min=1"`
	// You could add Offset int32 `form:"offset" binding:"min=0"`
}

// listVenues handles fetching a list of all venues.
// GET /venues
func (server *Server) ListVenues(ctx *gin.Context) {
	var req listVenuesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	venues, err := server.store.Listvenues(ctx) // Assuming ListVenues takes a context and no other params for now
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, venues)
}

// updateVenueRequest defines the request body for updating a venue.
type updateVenueRequest struct {
	ImageLinks      []string  `json:"image_links"`
	Name            string    `json:"name" binding:"required"`
	Type            string    `json:"type"`
	Description     string    `json:"description"`
	Location        string    `json:"location" binding:"required"`
	Dimension       string    `json:"dimension"`
	Capacity        int       `json:"capacity"`
	Facilities      []string  `json:"facilities"`
	HasAccomodation bool      `json:"has_accomodation"`
	RoomType        string    `json:"room_type"`
	NoOfRooms       int       `json:"no_of_rooms"`
	Sleeps          string    `json:"sleeps"`
	BedType         string    `json:"bed_type"`
	Rent            int       `json:"rent"`
	OwnedBy         uuid.UUID `json:"owned_by" binding:"required"`
	IsAvailable     bool      `json:"is_available"`
	OpensAt         time.Time `json:"opens_at"`
	ClosesAt        time.Time `json:"closes_at"`
	RentalDays      string    `json:"rental_days"`
	BookingPrice    int       `json:"booking_price"`
}

// updateVenueURI defines the URI parameter for updating a venue by ID.
type updateVenueURI struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// UpdateVenue handles updating an existing venue.
// PUT /venues/:id
func (server *Server) UpdateVenue(ctx *gin.Context) {
	var uri updateVenueURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	uuidID, err := uuid.Parse(uri.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid venue ID format: %w", err)))
		return
	}

	var req updateVenueRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateVenueParams{
		ID:              uuidID,
		Name:            req.Name,
		Column3:         req.ImageLinks,
		Type:            newNullString(req.Type),
		Description:     newNullString(req.Description),
		Location:        req.Location,
		Dimension:       newNullString(req.Dimension),
		Capacity:        newNullInt32(req.Capacity),
		Column9:         req.Facilities,
		HasAccomodation: sql.NullBool{Bool: req.HasAccomodation, Valid: true},
		RoomType:        newNullString(req.RoomType),
		NoOfRooms:       newNullInt32(req.NoOfRooms),
		Sleeps:          newNullString(req.Sleeps),
		BedType:         newNullString(req.BedType),
		Rent:            newNullInt32(req.Rent),
		OwnedBy:         newNullUUID(req.OwnedBy),
		IsAvailable:     sql.NullBool{Bool: req.IsAvailable, Valid: true},
		OpensAt:         newNullTime(req.OpensAt),
		ClosesAt:        newNullTime(req.ClosesAt),
		RentalDays:      newNullString(req.RentalDays),
		BookingPrice:    newNullInt32(req.BookingPrice),
	}

	err = server.store.UpdateVenue(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("venue not found for update")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Venue updated successfully",
		"id":      uuidID.String(),
	})
}

// deleteVenueRequest defines the URI parameter for deleting a venue by ID.
type deleteVenueRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// DeleteVenue handles deleting a venue by ID.
// DELETE /venues/:id
func (server *Server) DeleteVenue(ctx *gin.Context) {
	var req deleteVenueRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	uuidID, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid venue ID format: %w", err)))
		return
	}

	err = server.store.DeleteVenue(ctx, uuidID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("venue not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Venue deleted successfully",
	})
}
