package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"log/slog"

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
	userID := getUserID(c, defaultEnvProvider, defaultParserFactory)
	if userID == nil {
		return
	}

	// Get user's tags with message counts
	tags, err := getUserTagsWithCounts(db, *userID)
	if err != nil {
		slog.Error("Database error", "user_id", *userID, "error", err)

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

type EnvProvider interface {
	GetBotToken() string
}
type prodEnvProvider struct{}

func (m *prodEnvProvider) GetBotToken() string {
	return getBotToken()
}

var defaultEnvProvider = &prodEnvProvider{}

func getUserID(c *gin.Context, p EnvProvider, factory ParserFactory) *int64 {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "Authorization header is required!",
		})
		return nil
	}

	// Extract user ID from Telegram Web App data
	userID, err := extractUserIDFromAuth(authHeader, p, factory)
	if err != nil {
		slog.Error("Authentication error", "error", err)
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "Invalid authentication data",
		})
		return nil
	}
	return &userID
}

func getTagID(c *gin.Context) *int64 {
	tagIDStr := c.Param("tagId")
	tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
	if err != nil {
		slog.Error("Invalid tagId parameter", "tag_id_str", tagIDStr, "error", err)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid tag ID format",
		})
		return nil
	}
	return &tagID
}

func getTagMessagesHandler(c *gin.Context, db *sql.DB) {
	// Get authorization header
	userID := getUserID(c, defaultEnvProvider, defaultParserFactory)
	if userID == nil {
		return
	}

	// Get and validate tagId parameter
	tagID := getTagID(c)
	if tagID == nil {
		return
	}

	// Get messages for the specified tag
	messages, err := getTagMessages(db, *userID, *tagID)
	if err != nil {
		printMessagesError(c, userID, tagID, err)
		return
	}

	slog.Info("Successfully retrieved messages",
		"message_count", len(messages),
		"user_id", *userID,
		"tag_id", *tagID)

	// Return successful response
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    messages,
	})
}

func printMessagesError(c *gin.Context, userID *int64, tagID *int64, err error) {
	slog.Error("Database error", "user_id", *userID, "tag_id", *tagID, "error", err)

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
