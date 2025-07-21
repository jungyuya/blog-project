package main

import (
	"context" // Lambda 핸들러에 context 파라미터 필요
	"fmt"     // 디버깅용으로 임시 추가

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// handler 함수는 APIGatewayProxyRequest를 받아 APIGatewayProxyResponse를 반환합니다.
// 이것이 API Gateway와의 통합에 필요한 표준 형식입니다.
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("GET /posts 요청이 Go Lambda에 도착했습니다!") // CloudWatch 로그 확인용

	// 성공적인 HTTP 200 응답을 반환합니다.
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       `{"message": "Hello from Go Lambda for GET /posts!"}`,
	}, nil
}

func main() {
	// lambda.Start는 handler 함수를 Lambda 런타임에 등록합니다.
	lambda.Start(handler)
}
