package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time" // UpdatedAt 필드 업데이트를 위해 필요

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
// ⭐ Author 필드가 여기에 추가되어 있어야 합니다. ⭐
type Post struct {
	PostID    string `json:"postId" dynamodbav:"postId"`
	Title     string `json:"title" dynamodbav:"title"`
	Content   string `json:"content" dynamodbav:"content"`
	Author    string `json:"author" dynamodbav:"author"` // ⭐ 이 필드가 있는지 확인
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

// handleRequest는 updatePost Lambda 함수의 메인 핸들러입니다.
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (apiResponse, error) {
	// API Gateway 요청에서 경로 파라미터 'postId' 값을 가져옵니다.
	postId := request.PathParameters["postId"]
	fmt.Printf("updatePost 함수 호출됨. PostID: %s 업데이트 요청.\n", postId)

	if postId == "" {
		fmt.Println("PostID가 요청 경로에 없습니다.")
		return apiGatewayResponse(400, `{"message": "PostID is required"}`), nil
	}

	// 1. 요청 본문 파싱 (업데이트할 데이터)
	// 클라이언트로부터 전송된 JSON 요청 본문을 Go의 임시 구조체에 매핑합니다.
	// ⭐ 여기에 Author 필드를 추가합니다. ⭐
	var updatedPostData struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Author  string `json:"author"` // ⭐ 이 줄을 추가합니다. ⭐
	}
	err := json.Unmarshal([]byte(request.Body), &updatedPostData)
	if err != nil {
		fmt.Printf("요청 본문 파싱 오류: %v\n", err)
		return apiGatewayResponse(400, `{"message": "Invalid request body"}`), nil
	}

	// 유효성 검사: 업데이트할 내용이 비어있는지 확인합니다.
	// ⭐ Author 유효성 검사 조건도 추가합니다. ⭐
	if updatedPostData.Title == "" || updatedPostData.Content == "" || updatedPostData.Author == "" {
		fmt.Println("업데이트할 제목, 내용 또는 작성자가 비어 있습니다.")
		return apiGatewayResponse(400, `{"message": "Title, content, and author must be provided for update"}`), nil
	}

	// 2. DynamoDB UpdateItem API 호출을 위한 Input 구성
	// `UpdateExpression`은 DynamoDB 항목을 어떻게 업데이트할지 정의하는 문자열입니다.
	// ⭐ #A (Author)와 :a (author value)를 추가합니다. ⭐
	updateExpression := "SET #T = :t, #C = :c, #A = :a, #U = :u"
	
	// `ExpressionAttributeNames`는 DynamoDB 예약어(예: "Content")와 충돌을 피하기 위해
	// 실제 속성 이름을 플레이스홀더에 매핑합니다.
	// ⭐ #A (Author)를 추가합니다. ⭐
	expressionAttributeNames := map[string]*string{
		"#T": aws.String("title"),
		"#C": aws.String("content"),
		"#A": aws.String("author"), // ⭐ 이 줄을 추가합니다. ⭐
		"#U": aws.String("updatedAt"),
	}

	// `ExpressionAttributeValues`는 UpdateExpression에서 사용될 실제 값들을 정의합니다.
	// ⭐ :a (author value)를 추가합니다. ⭐
	now := time.Now().UTC().Format(time.RFC3339)
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":t": {S: aws.String(updatedPostData.Title)},
		":c": {S: aws.String(updatedPostData.Content)},
		":a": {S: aws.String(updatedPostData.Author)}, // ⭐ 이 줄을 추가합니다. ⭐
		":u": {S: aws.String(now)},
	}

	// `UpdateItemInput`은 DynamoDB에 항목을 업데이트하기 위한 요청 파라미터를 정의합니다.
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName), // 업데이트할 테이블 이름
		Key: map[string]*dynamodb.AttributeValue{ // 업데이트할 항목의 기본 키
			"postId": {
				S: aws.String(postId),
			},
		},
		UpdateExpression:         aws.String(updateExpression),          // 업데이트 표현식
		ExpressionAttributeNames: expressionAttributeNames,              // 속성 이름 매핑
		ExpressionAttributeValues: expressionAttributeValues,             // 속성 값 매핑
		ReturnValues:             aws.String("ALL_NEW"),                 // 업데이트된 항목의 모든 새 속성을 반환하도록 요청
	}

	result, err := dynamoClient.UpdateItem(input)
	if err != nil {
		// DynamoDB UpdateItem 중 오류 발생 시, 500 Internal Server Error 응답을 반환합니다.
		fmt.Printf("DynamoDB UpdateItem 오류: %v\n", err)
		return apiGatewayResponse(500, `{"message": "Failed to update post"}`), nil
	}

	// 3. 업데이트 결과 확인 및 언마샬링
	// `result.Attributes`는 업데이트된 항목의 모든 새 속성을 포함하는 AttributeValue 맵입니다.
	if result.Attributes == nil || len(result.Attributes) == 0 {
		// 업데이트할 항목을 찾지 못했거나, 업데이트 후 반환된 속성이 없는 경우 (매우 드물게 발생)
		fmt.Printf("PostID: %s 에 해당하는 게시물을 찾을 수 없거나 업데이트 실패.\n", postId)
		return apiGatewayResponse(404, `{"message": "Post not found or update failed"}`), nil
	}

	var updatedPost Post
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &updatedPost)
	if err != nil {
		// 언마샬링 오류 발생 시, 500 Internal Server Error 응답을 반환합니다.
		fmt.Printf("DynamoDB 항목 언마샬링 오류: %v\n", err)
		return apiGatewayResponse(500, `{"message": "Failed to process updated post data"}`), nil
	}

	fmt.Printf("PostID: %s 게시물 성공적으로 업데이트됨.\n", postId)

	// 4. 성공 응답 반환
	// 업데이트된 Post 객체를 JSON 문자열로 마샬링하여 응답 본문으로 전송합니다.
	responseBody, err := json.Marshal(updatedPost)
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
