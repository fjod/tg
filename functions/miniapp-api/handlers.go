package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

		api.GET("/user/tags/:tagId/messages", func(c *gin.Context) {
			log.Printf("=== ROUTE MATCHED: /user/tags/:tagId/messages ===")
			getTagMessagesHandler(c, db)
		})
		api.OPTIONS("/user/tags/:tagId/messages", optionsHandler)
	}

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Allow all yandexcloud.net and website.yandexcloud.net domains
		if origin != "" && (containsYandexDomain(origin)) {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "false")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func containsYandexDomain(origin string) bool {
	return origin != "" && (
	// Allow both API gateway and Object Storage domains
	origin == "https://d5di1npf8thkd9m534rv.8wihnuyr.apigw.yandexcloud.net" ||
		origin == "https://tg-bot-storage-fjod.website.yandexcloud.net" ||
		// Allow any yandexcloud.net subdomain for flexibility
		(len(origin) > 16 && origin[:8] == "https://" &&
			(origin[len(origin)-16:] == ".yandexcloud.net" ||
				origin[len(origin)-24:] == ".website.yandexcloud.net")))
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
			Error:   "Authorization header is required!",
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

func getTagMessagesHandler(c *gin.Context, db *sql.DB) {
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

	// Get and validate tagId parameter
	tagIdStr := c.Param("tagId")
	log.Printf("=== TAG ID PARSING DEBUG ===")
	log.Printf("Raw tagIdStr from c.Param('tagId'): '%s'", tagIdStr)
	log.Printf("Length of tagIdStr: %d", len(tagIdStr))
	log.Printf("Request URL: %s", c.Request.URL.String())
	log.Printf("Request Path: %s", c.Request.URL.Path)

	// Also try getting all params to see what Gin has
	allParams := c.Params
	log.Printf("All Gin params: %+v", allParams)

	tagId, err := strconv.ParseInt(tagIdStr, 10, 64)
	if err != nil {
		log.Printf("FAILED to parse tagId: '%s', error: %v", tagIdStr, err)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Invalid tag ID format: '%s'", tagIdStr),
		})
		return
	}

	log.Printf("Successfully parsed tagId: %d", tagId)
	log.Printf("=== END TAG ID PARSING DEBUG ===")

	// Get messages for the specified tag
	messages, err := getTagMessages(db, userID, tagId)
	if err != nil {
		log.Printf("Database error for user %d, tag %d: %v", userID, tagId, err)

		// Check if it's a not found/access denied error
		if err.Error() == "tag not found or access denied" {
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Error:   "Tag not found or you don't have access to it",
			})
			return
		}

		// General database error
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to fetch messages for tag",
		})
		return
	}

	log.Printf("Successfully retrieved %d messages for user %d, tag %d", len(messages), userID, tagId)

	// Return successful response
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    messages,
	})
}
