package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func setupRoutes(db *sql.DB) *gin.Engine {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Add logging middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Add CORS middleware
	r.Use(corsMiddleware())

	// API routes
	api := r.Group("/api")
	{
		// Health check endpoint (no auth required)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, APIResponse{
				Success: true,
				Data:    map[string]string{"status": "healthy", "timestamp": "2025-01-15T12:00:00Z"},
			})
		})
		api.OPTIONS("/health", optionsHandler)

		api.GET("/user/tags", func(c *gin.Context) {
			getUserTagsHandler(c, db)
		})
		api.OPTIONS("/user/tags", optionsHandler)
	}

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func optionsHandler(c *gin.Context) {
	// OPTIONS requests are handled by CORS middleware
	// Just return 200 OK status
	c.Status(http.StatusOK)
}

func getUserTagsHandler(c *gin.Context, db *sql.DB) {
	// Get authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "Authorization header is required",
		})
		return
	}

	// Extract user ID from Telegram Web App data
	userID, err := extractUserIDFromAuth(authHeader)
	if err != nil {
		log.Printf("Authentication error: %v", err)
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "Invalid authentication data",
		})
		return
	}

	// Get user's tags with message counts
	tags, err := getUserTagsWithCounts(db, userID)
	if err != nil {
		log.Printf("Database error for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to fetch user tags",
		})
		return
	}

	// Return successful response
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    tags,
	})
}
