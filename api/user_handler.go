package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync" // Import sync for mutex
	"time" // Import time for JWT expiration

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5" // Import the JWT library
	"github.com/google/uuid"
	db "github.com/tedobanks/tabularasa_backend/db/sqlc" // Assuming this path is correct
	"golang.org/x/crypto/bcrypt"
)

// Define JWT secret and token duration (these should ideally come from environment variables or a config file)
const (
	// jwtSecret is the secret key used to sign JWT tokens.
	// In a production environment, this should be a strong, randomly generated key
	// loaded from a secure source (e.g., environment variable, secret management service).
	jwtSecret = "your_super_secret_jwt_key_that_is_at_least_32_bytes_long" // CHANGE THIS IN PRODUCTION!
	// tokenDuration specifies how long the JWT token is valid for.
	tokenDuration = time.Hour * 24 // Tokens expire after 24 hours
)

// tokenBlacklist is an in-memory map to store blacklisted JWTs (using their JTI).
// In a production environment, this should be a persistent store like Redis or a database.
var tokenBlacklist = make(map[string]time.Time)

// blacklistMutex protects concurrent access to the tokenBlacklist map.
var blacklistMutex sync.Mutex

// Helper function to create sql.NullString from a string
func newNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// Function to hash a password
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// Function to check if a password matches a hash
func checkPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// makeToken generates a new JWT token for a given user ID.
// The token contains the user's ID, a unique JTI, and an expiration time.
func makeToken(userID uuid.UUID) (string, error) {
	// Generate a unique ID for the token (JTI - JWT ID)
	tokenID := uuid.New().String()

	// Define the claims for the JWT token.
	// jwt.RegisteredClaims includes standard claims like Issuer, Subject, Audience, ExpirationTime, etc.
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)), // Token expiration time
		Subject:   userID.String(),                                   // Subject of the token (user ID)
		ID:        tokenID,                                           // Unique ID for the token (JTI)
	}

	// Create a new token with the specified signing method and claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key.
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// createUserRequest defines the request body for creating a new user (registration).
// The 'Roles' field has been removed as it will now be defaulted to "Personal Profile" internally.
type createUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// createUserResponse defines the response structure for successful user creation and login.
type createUserResponse struct {
	User    db.Users    `json:"user"`
	Profile db.Profiles `json:"profile"`
	Token   string      `json:"token"`
}

// CreateUser handles the creation of a new user (registration) and automatically creates a profile for them.
// It also generates and returns a JWT token upon successful registration.
// The profile's 'Roles' will now default to "Personal Profile".
// POST /register
func (server *Server) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Prepare arguments for user creation
	userArg := db.CreateUserParams{
		Email:     req.Email,
		Password:  newNullString(hashedPassword),
		Firstname: newNullString(req.Firstname),
		Lastname:  newNullString(req.Lastname),
	}

	// Create the user in the database
	user, err := server.store.CreateUser(ctx, userArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Prepare arguments for profile creation using the newly created user's ID.
	// The 'Roles' field is now hardcoded to "Personal Profile".
	profileArg := db.CreateProfileParams{
		ID:    user.ID,
		Roles: "Personal Profile", // Default role set here
	}

	// Create the profile in the database
	profile, err := server.store.CreateProfile(ctx, profileArg)
	if err != nil {
		// If profile creation fails, consider rolling back the user creation (requires database transactions).
		// For simplicity, we'll just return an error here.
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create profile for user %s: %w", user.ID, err)))
		return
	}

	// Generate JWT token for the newly created user
	token, err := makeToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to generate token: %w", err)))
		return
	}

	// Respond with user, profile, and the generated token
	rsp := createUserResponse{
		User:    user,
		Profile: profile,
		Token:   token,
	}
	ctx.JSON(http.StatusOK, rsp)
}

// loginUserRequest defines the request body for user login.
type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginUser handles user authentication and returns a JWT token upon successful login.
// POST /login
func (server *Server) LoginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Retrieve user by email
	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("user with this email not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Check if the provided password matches the hashed password in the database
	if user.Password.Valid {
		err = checkPassword(req.Password, user.Password.String)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("incorrect password")))
			return
		}
	} else {
		// Handle cases where password might be null (e.g., if user was created via social login)
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("password not set for this user")))
		return
	}

	// Retrieve the user's profile
	profile, err := server.store.GetProfile(ctx, user.ID) // Assuming you have a GetProfile method in your store
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("profile not found for user")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Generate JWT token for the authenticated user
	token, err := makeToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to generate token: %w", err)))
		return
	}

	// Respond with user, profile, and the generated token
	rsp := createUserResponse{ // Reusing createUserResponse as it has the same structure
		User:    user,
		Profile: profile,
		Token:   token,
	}
	ctx.JSON(http.StatusOK, rsp)
}

// LogoutUser handles user logout by blacklisting the current token.
// POST /logout (requires authentication)
func (server *Server) LogoutUser(ctx *gin.Context) {
	// Get the token string from the Authorization header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("authorization header is missing")))
		return
	}
	tokenString := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	} else {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid authorization header format")))
		return
	}

	// Parse the token to get its claims, specifically the JTI and expiration
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("invalid token for logout: %w", err)))
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("invalid token claims for logout")))
		return
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("token does not contain JTI claim")))
		return
	}

	// Extract expiration time from claims to know when to remove from blacklist
	exp, ok := claims["exp"].(float64) // JWT 'exp' is typically a numeric date (Unix timestamp)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("token does not contain expiration claim")))
		return
	}
	expirationTime := time.Unix(int64(exp), 0)

	// Add the JTI to the blacklist with its expiration time
	blacklistMutex.Lock()
	tokenBlacklist[jti] = expirationTime
	blacklistMutex.Unlock()

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// AuthMiddleware is a Gin middleware to authenticate requests using JWT.
// It extracts the token from the Authorization header, verifies it,
// and sets the authenticated user's ID in the Gin context.
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("authorization header is required")))
			return
		}

		// Expected format: "Bearer <token>"
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("invalid authorization header format")))
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil // Return the secret key for verification
		})

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("invalid token: %w", err)))
			return
		}

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("invalid token claims")))
			return
		}

		// Check if the token's JTI is in the blacklist
		jti, jtiOk := claims["jti"].(string)
		if jtiOk {
			blacklistMutex.Lock()
			blacklistedExp, found := tokenBlacklist[jti]
			blacklistMutex.Unlock()

			// If found and the token has not yet expired naturally, it's revoked
			if found && time.Now().Before(blacklistedExp) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("token has been revoked")))
				return
			}
		}

		// Extract user ID from claims and set it in the Gin context
		userIDStr, ok := claims["sub"].(string) // "sub" is the standard claim for Subject
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("user ID not found in token claims")))
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("invalid user ID format in token claims")))
			return
		}

		ctx.Set("userID", userID) // Set the user ID in the context for subsequent handlers
		ctx.Next()                // Proceed to the next handler
	}
}

// getUserRequest defines the URI parameter for getting a user by ID.
type getUserRequest struct {
	ID string `uri:"id" binding:"required,uuid"` // Use 'uuid' binding for UUID string format
}

// getUser handles fetching a single user by ID.
// GET /users/:id
func (server *Server) GetUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	uuidID, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid user ID format: %w", err)))
		return
	}

	user, err := server.store.GetUser(ctx, uuidID) // <--- Pass uuid.UUID
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// getUserByEmailRequest defines the query parameter for getting a user by email.
type getUserByEmailRequest struct {
	Email string `form:"email" binding:"required,email"`
}

// getUserByEmail handles fetching a single user by email.
// GET /users/by-email?email=...
func (server *Server) GetUserByEmail(ctx *gin.Context) {
	var req getUserByEmailRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// listUsersRequest defines optional query parameters for listing users (e.g., pagination).
type listUsersRequest struct {
	// You could add Limit int64 `form:"limit" binding:"min=1"`
	// You could add Offset int64 `form:"offset" binding:"min=0"`
}

// listUsers handles fetching a list of all users.
// GET /users
func (server *Server) ListUsers(ctx *gin.Context) {
	var req listUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	users, err := server.store.ListUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// updateuserRequest defines the request body for updating a user.
type updateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// updateUserURI defines the URI parameter for updating a user by ID.
type updateUserURI struct {
	ID string `uri:"id" binding:"required,uuid"` // Use 'uuid' binding for UUID string format
}

// updateUser handles updating an existing user.
// PUT /users/:id
func (server *Server) UpdateUser(ctx *gin.Context) {
	var uri updateUserURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	uuidID, err := uuid.Parse(uri.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid user ID format: %w", err)))
		return
	}

	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		ID:        uuidID, // <--- Pass uuid.UUID
		Email:     req.Email,
		Password:  newNullString(hashedPassword),
		Firstname: newNullString(req.Firstname),
		Lastname:  newNullString(req.Lastname),
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("user not found for update")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// deleteUserRequest defines the URI parameter for deleting a user by ID.
type deleteUserRequest struct {
	ID string `uri:"id" binding:"required,uuid"` // Use 'uuid' binding for UUID string format
}

// deleteUser handles deleting a user by ID.
// DELETE /users/:id
func (server *Server) DeleteUser(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	uuidID, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid user ID format: %w", err)))
		return
	}

	err = server.store.DeleteUser(ctx, uuidID)
	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
