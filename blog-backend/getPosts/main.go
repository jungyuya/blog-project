package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	// "github.com/google/uuid" // getPosts에서는 UUID 생성 필요 없음
)

// DynamoDB 클라이언트 및 테이블 이름 전역 변수 (createPost와 동일)
var (
	dynamoClient dynamodbiface.DynamoDBAPI
	tableName    string
)

func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamoClient = dynamodb.New(sess)
	tableName = os.Getenv("DYNAMODB_TABLE")
	if tableName == "" {
		fmt.Println("경고: 환경 변수 DYNAMODB_TABLE이 설정되지 않았습니다.")
	}
}

// Post 구조체 (createPost와 동일하게 정의)
type Post struct {
	PostID    string `json:"postId" dynamodbav:"postId"`
	Title     string `json:"title" dynamodbav:"title"`
	Content   string `json:"content" dynamodbav:"content"`
	CreatedAt string `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt string `json:"updatedAt" dynamodbav:"updatedAt"`
}

// apiResponse 구조체 (createPost와 동일)
type apiResponse struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// apiGatewayResponse 헬퍼 함수 (createPost와 동일)
func apiGatewayResponse(statusCode int, body string) apiResponse {
	return apiResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		},
		Body: body,
	}
}

// handleRequest는 getPosts Lambda 함수의 메인 핸들러입니다.
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (apiResponse, error) {
	fmt.Println("getPosts 함수 호출됨. 모든 게시물 조회 요청.")

	// 1. DynamoDB Scan API 호출
	// ScanInput은 DynamoDB 테이블의 모든 항목을 가져오기 위한 요청 파라미터를 정의합니다.
	// 주의: 대량의 데이터셋에서는 Scan 작업이 성능 및 비용 측면에서 비효율적일 수 있습니다.
	// 실제 운영 환경에서는 Query 또는 다른 최적화된 방법을 고려해야 합니다.
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName), // 조회할 테이블 이름
	}

	result, err := dynamoClient.Scan(input)
	if err != nil {
		// DynamoDB Scan 중 오류 발생 시, 500 Internal Server Error 응답을 반환합니다.
		fmt.Printf("DynamoDB Scan 오류: %v\n", err)
		return apiGatewayResponse(500, `{"message": "Failed to retrieve posts"}`), nil
	}

	// 2. DynamoDB 항목 목록을 Go 슬라이스(Slice)로 언마샬링
	// `dynamodbattribute.UnmarshalListOfMaps` 함수는 DynamoDB의 AttributeValue 맵 목록을
	// Go 구조체 슬라이스([]Post)로 변환합니다.
	var posts []Post
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &posts)
	if err != nil {
		// 언마샬링 오류 발생 시, 500 Internal Server Error 응답을 반환합니다.
		fmt.Printf("DynamoDB 항목 목록 언마샬링 오류: %v\n", err)
		return apiGatewayResponse(500, `{"message": "Failed to process retrieved posts"}`), nil
	}

	fmt.Printf("총 %d개의 게시물 조회됨.\n", len(posts))

	// 3. 성공 응답 반환
	// 조회된 Post 객체 슬라이스를 JSON 문자열로 마샬링하여 응답 본문으로 전송합니다.
	responseBody, err := json.Marshal(posts)
	if err != nil {
		fmt.Printf("응답 본문 마샬링 오류: %v\n", err)
		return apiGatewayResponse(500, `{"message": "Failed to marshal response"}`), nil
	}

	// 200 OK HTTP 상태 코드와 함께 성공 응답을 반환합니다.
	return apiGatewayResponse(200, string(responseBody)), nil
}

func main() {
	lambda.Start(handleRequest)
}