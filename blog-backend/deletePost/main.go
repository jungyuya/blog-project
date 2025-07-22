package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

// handleRequest는 deletePost Lambda 함수의 메인 핸들러입니다.
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (apiResponse, error) {
	// API Gateway 요청에서 경로 파라미터 'postId' 값을 가져옵니다.
	postId := request.PathParameters["postId"]
	fmt.Printf("deletePost 함수 호출됨. PostID: %s 삭제 요청.\n", postId)

	if postId == "" {
		fmt.Println("PostID가 요청 경로에 없습니다.")
		return apiGatewayResponse(400, `{"message": "PostID is required"}`), nil
	}

	// 1. DynamoDB DeleteItem API 호출을 위한 Input 구성
	// `DeleteItemInput`은 DynamoDB 테이블에서 특정 항목을 삭제하기 위한 요청 파라미터를 정의합니다.
	// `Key`는 삭제하려는 항목의 기본 키(Primary Key)를 나타내는 AttributeValue 맵입니다.
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName), // 삭제할 테이블 이름
		Key: map[string]*dynamodb.AttributeValue{
			"postId": { // 파티션 키 이름
				S: aws.String(postId), // 삭제하려는 postId 값
			},
		},
		// ReturnValues: aws.String("ALL_OLD"), // 삭제된 항목을 반환하려면 이 옵션 사용 (선택 사항)
	}

	_, err := dynamoClient.DeleteItem(input)
	if err != nil {
		// DynamoDB DeleteItem 중 오류 발생 시, 500 Internal Server Error 응답을 반환합니다.
		fmt.Printf("DynamoDB DeleteItem 오류: %v\n", err)
		return apiGatewayResponse(500, `{"message": "Failed to delete post"}`), nil
	}

	fmt.Printf("PostID: %s 게시물 성공적으로 삭제됨.\n", postId)

	// 2. 성공 응답 반환
	// 삭제 작업은 보통 반환할 내용이 없으므로 204 No Content 상태 코드를 사용합니다.
	// API Gateway는 204 No Content 응답을 받으면 본문을 보내지 않습니다.
	return apiGatewayResponse(204, ""), nil // 204 응답은 본문이 없습니다.
}

func main() {
	lambda.Start(handleRequest)
}