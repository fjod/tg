package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var db *sql.DB

func getBotToken() string {
	return os.Getenv("TELEGRAM_BOT_TOKEN")
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
		log.Printf("Failed to convert Lambda request: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"success": false, "error": "Invalid request format"}`,
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

	// Convert to Lambda response
	return events.APIGatewayProxyResponse{
		StatusCode: recorder.statusCode,
		Body:       recorder.body,
		Headers:    recorder.headers,
	}, nil
}

func convertLambdaRequest(request events.APIGatewayProxyRequest) (*http.Request, error) {
	// Create HTTP request from Lambda request
	req, err := http.NewRequest(request.HTTPMethod, request.Path, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	// Add query parameters
	q := req.URL.Query()
	for key, value := range request.QueryStringParameters {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

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
