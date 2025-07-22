package main

import (
    "fmt"
    "context"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // 이 곳에 createPost 로직을 추가할 것입니다.
    fmt.Println("POST /posts 요청이 Go Lambda에 도착했습니다!")
    return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Headers:    map[string]string{"Content-Type": "application/json"},
        Body:       `{"message": "Hello from Go Lambda for POST /posts!"}`,
    }, nil
}

func main() {
    lambda.Start(handler)
}