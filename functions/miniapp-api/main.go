package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var db *sql.DB

func getBotToken() string {
	return os.Getenv("TELEGRAM_BOT_TOKEN")
}

func containsPattern(origin, pattern string) bool {
	return strings.Contains(origin, pattern)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Log incoming request details
	log.Printf("=== LAMBDA REQUEST RECEIVED ===")
	log.Printf("HTTP Method: %s", request.HTTPMethod)
	log.Printf("Path: %s", request.Path)
	log.Printf("Resource: %s", request.Resource)
	log.Printf("Stage: %s", request.RequestContext.Stage)
	log.Printf("Headers: %+v", request.Headers)
	log.Printf("Query Params: %+v", request.QueryStringParameters)
	log.Printf("Body: %s", request.Body)
	log.Printf("==============================")

	// Initialize database connection if not already done
	if db == nil {
		var err error
		db, err = initDB()
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"success": false, "error": "Database connection failed"}`,
				Headers: map[string]string{
					"Content-Type":                "application/json",
					"Access-Control-Allow-Origin": "*",
				},
			}, nil
		}
	}

	// Create Gin router
	router := setupRoutes(db)

	// Convert Lambda request to HTTP request
	req, err := convertLambdaRequest(request)
	if err != nil {
		log.Printf("Failed to convert Lambda request : %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"success": false, "error": " Invalid request format"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Create response recorder
	recorder := &ResponseRecorder{
		headers: make(map[string]string),
	}

	// Process the request
	router.ServeHTTP(recorder, req)

	// Ensure we always have a status code
	if recorder.statusCode == 0 {
		recorder.statusCode = 200
	}

	// Ensure CORS headers are always present
	if recorder.headers == nil {
		recorder.headers = make(map[string]string)
	}

	// Set CORS headers to allow both domain patterns
	origin := request.Headers["origin"]
	if origin == "" {
		origin = request.Headers["Origin"]
	}

	// Allow both yandexcloud.net and website.yandexcloud.net domains
	if origin != "" && (containsPattern(origin, "yandexcloud.net") ||
		containsPattern(origin, "website.yandexcloud.net")) {
		recorder.headers["Access-Control-Allow-Origin"] = origin
	} else {
		recorder.headers["Access-Control-Allow-Origin"] = "*"
	}

	recorder.headers["Access-Control-Allow-Methods"] = "GET, POST, PUT, DELETE, OPTIONS"
	recorder.headers["Access-Control-Allow-Headers"] = "Origin, Content-Type, Authorization"
	recorder.headers["Access-Control-Allow-Credentials"] = "false"

	log.Printf("Returning response - Status: %d, Body length: %d, Headers: %+v",
		recorder.statusCode, len(recorder.body), recorder.headers)

	// Convert to Lambda response
	return events.APIGatewayProxyResponse{
		StatusCode: recorder.statusCode,
		Body:       recorder.body,
		Headers:    recorder.headers,
	}, nil
}

func convertLambdaRequest(request events.APIGatewayProxyRequest) (*http.Request, error) {
	// Determine the correct path to use
	path := request.Path
	if path == "" {
		path = request.Resource
	}

	// Replace path parameters in the path
	// API Gateway gives us path parameters like {tagId} in PathParameters
	if len(request.PathParameters) > 0 {
		for key, value := range request.PathParameters {
			placeholder := "{" + key + "}"
			path = strings.Replace(path, placeholder, value, -1)
		}
	}

	log.Printf("=== LAMBDA REQUEST CONVERSION DEBUG ===")
	log.Printf("Original Path: '%s'", request.Path)
	log.Printf("Resource: '%s'", request.Resource)
	log.Printf("PathParameters: %+v", request.PathParameters)
	log.Printf("Final Path after substitution: '%s'", path)
	log.Printf("HTTP Method: %s", request.HTTPMethod)
	log.Printf("=== END LAMBDA REQUEST CONVERSION DEBUG ===")

	// Create HTTP request from Lambda request
	req, err := http.NewRequest(request.HTTPMethod, path, nil)
	if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
		return nil, err
	}

	// Add headers
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	log.Printf("Added %d headers to request", len(request.Headers))

	// Add query parameters
	q := req.URL.Query()
	for key, value := range request.QueryStringParameters {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	log.Printf("Final request URL: %s", req.URL.String())

	return req, nil
}

type ResponseRecorder struct {
	statusCode int
	body       string
	headers    map[string]string
}

func (r *ResponseRecorder) Header() http.Header {
	h := make(http.Header)
	for key, value := range r.headers {
		h.Set(key, value)
	}
	return h
}

func (r *ResponseRecorder) Write(data []byte) (int, error) {
	r.body = string(data)
	return len(data), nil
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}

func main() {
	lambda.Start(Handler)
}
