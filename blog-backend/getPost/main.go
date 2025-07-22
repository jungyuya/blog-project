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
)

// DynamoDB 클라이언트 및 테이블 이름 전역 변수 (다른 함수들과 동일)
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

// Post 구조체 (다른 함수들과 동일하게 정의)
type Post struct {
	PostID    string `json:"postId" dynamodbav:"postId"`
	Title     string `json:"title" dynamodbav:"title"`
	Content   string `json:"content" dynamodbav:"content"`
	CreatedAt string `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt string `json:"updatedAt" dynamodbav:"updatedAt"`
}

// apiResponse 구조체 (다른 함수들과 동일)
type apiResponse struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// apiGatewayResponse 헬퍼 함수 (다른 함수들과 동일)
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

// handleRequest는 getPost Lambda 함수의 메인 핸들러입니다.
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (apiResponse, error) {
	// API Gateway 요청에서 경로 파라미터(pathParameter) 'postId' 값을 가져옵니다.
	postId := request.PathParameters["postId"]
	fmt.Printf("getPost 함수 호출됨. PostID: %s 조회 요청.\n", postId)

	if postId == "" {
		fmt.Println("PostID가 요청 경로에 없습니다.")
		return apiGatewayResponse(400, `{"message": "PostID is required"}`), nil
	}

	// 1. DynamoDB GetItem API 호출을 위한 Input 구성
	// `GetItemInput`은 DynamoDB 테이블에서 특정 항목을 가져오기 위한 요청 파라미터를 정의합니다.
	// `Key`는 조회하려는 항목의 기본 키(Primary Key)를 나타내는 AttributeValue 맵입니다.
	// 여기서는 `postId`가 파티션 키로 사용됩니다.
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName), // 조회할 테이블 이름
		Key: map[string]*dynamodb.AttributeValue{
			"postId": { // 파티션 키 이름 (DynamoDB 테이블 스키마에 따라 다름)
				S: aws.String(postId), // 조회하려는 postId 값 (문자열 타입)
			},
		},
	}

	result, err := dynamoClient.GetItem(input)
	if err != nil {
		// DynamoDB GetItem 중 오류 발생 시, 500 Internal Server Error 응답을 반환합니다.
		fmt.Printf("DynamoDB GetItem 오류: %v\n", err)
		return apiGatewayResponse(500, `{"message": "Failed to retrieve post"}`), nil
	}

	// 2. 조회 결과 확인 및 언마샬링
	// `result.Item`은 조회된 단일 항목의 AttributeValue 맵입니다.
	// 만약 항목이 존재하지 않으면 `result.Item`은 `nil` 또는 비어있는 맵이 됩니다.
	if result.Item == nil || len(result.Item) == 0 {
		fmt.Printf("PostID: %s 에 해당하는 게시물을 찾을 수 없습니다.\n", postId)
		return apiGatewayResponse(404, `{"message": "Post not found"}`), nil
	}

	// `dynamodbattribute.UnmarshalMap` 함수는 DynamoDB의 AttributeValue 맵을
	// Go 구조체(Post)로 변환합니다.
	var post Post
	err = dynamodbattribute.UnmarshalMap(result.Item, &post)
	if err != nil {
		// 언마샬링 오류 발생 시, 500 Internal Server Error 응답을 반환합니다.
		fmt.Printf("DynamoDB 항목 언마샬링 오류: %v\n", err)
		return apiGatewayResponse(500, `{"message": "Failed to process retrieved post"}`), nil
	}

	fmt.Printf("PostID: %s 게시물 성공적으로 조회됨.\n", postId)

	// 3. 성공 응답 반환
	// 조회된 Post 객체를 JSON 문자열로 마샬링하여 응답 본문으로 전송합니다.
	responseBody, err := json.Marshal(post)
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