// utils/cors.go
package utils

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

var allowedOrigins = []string{
	"http://localhost:3000",
	"https://blog.jungyu.store",
}

// ValidateOrigin checks the request's Origin header against allowedOrigins.
// Returns the origin if allowed, otherwise returns an empty string.
func ValidateOrigin(origin string) string {
	for _, o := range allowedOrigins {
		if strings.EqualFold(o, origin) {
			return origin
		}
	}
	fmt.Printf("❌ 허용되지 않은 Origin 요청: %s\n", origin)
	return ""
}

// NewCORSResponse generates an APIGatewayProxyResponse with CORS headers set.
func NewCORSResponse(statusCode int, body string, origin string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  origin,
			"Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		},
		Body: body,
	}
}
