package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
)

// DynamoDB 클라이언트를 전역으로 초기화합니다.
// 이는 Lambda 함수의 '콜드 스타트(Cold Start)' 성능을 최적화하기 위한 중요한 기법입니다.
// Lambda 함수가 처음 호출될 때 한 번만 초기화되고, 이후 재사용되기 때문입니다.
var (
	dynamoClient dynamodbiface.DynamoDBAPI // DynamoDB API 클라이언트 인터페이스
	tableName    string                    // DynamoDB 테이블 이름을 저장할 변수
)

// init() 함수는 main() 함수가 실행되기 전에 Go 런타임에 의해 자동으로 호출됩니다.
// 주로 초기화 작업을 수행하는 데 사용됩니다.
func init() {
	// AWS 세션 초기화: AWS 서비스와 상호작용하기 위한 기본 연결 설정입니다.
	// SharedConfigEnable 옵션은 ~/.aws/config 또는 환경 변수 등 공유 설정 파일을 사용하도록 합니다.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// DynamoDB 서비스 클라이언트 생성: 초기화된 세션을 바탕으로 DynamoDB API를 호출할 클라이언트를 만듭니다.
	dynamoClient = dynamodb.New(sess)

	// 환경 변수에서 DynamoDB 테이블 이름 가져오기:
	// Serverless Framework의 serverless.yml에 정의된 DYNAMODB_TABLE 환경 변수 값을 읽어옵니다.
	// 이는 환경(개발/스테이징/운영)에 따라 다른 테이블 이름을 유연하게 사용할 수 있도록 합니다.
	tableName = os.Getenv("DYNAMODB_TABLE")
	if tableName == "" {
		// 환경 변수가 설정되지 않았을 경우, 개발 중이거나 설정 오류일 수 있으므로 경고 메시지를 출력합니다.
		fmt.Println("경고: 환경 변수 DYNAMODB_TABLE이 설정되지 않았습니다.")
		// 실제 운영 환경에서는 이 경우 패닉(panic)을 발생시키거나 오류를 반환하여 함수 실행을 중단해야 합니다.
	}
}

// Post 구조체는 블로그 게시물의 데이터 모델을 정의합니다.
// `json:"..."` 태그는 이 구조체가 JSON 데이터를 마샬링(Go 객체를 JSON으로 변환)하거나
// 언마샬링(JSON을 Go 객체로 변환)할 때 사용될 필드 이름을 지정합니다.
// `dynamodbav:"..."` 태그는 이 구조체가 DynamoDB 항목(AttributeValue)으로 변환될 때 사용될 속성 이름을 지정합니다.
type Post struct {
	PostID    string `json:"postId" dynamodbav:"postId"`      // 게시물 고유 ID
	Title     string `json:"title" dynamodbav:"title"`        // 게시물 제목
	Content   string `json:"content" dynamodbav:"content"`    // 게시물 내용
	Author    string `json:"author" dynamodbav:"author"`      // 게시물 작성자 ⭐ 이 필드는 이미 추가되어 있습니다. ⭐
	CreatedAt string `json:"createdAt" dynamodbav:"createdAt"`// 게시물 생성 시각 (ISO 8601 형식)
	UpdatedAt string `json:"updatedAt" dynamodbav:"updatedAt"`// 게시물 마지막 업데이트 시각 (ISO 8601 형식)
}

// apiResponse 구조체는 API Gateway에 반환할 HTTP 응답의 표준 형식을 정의합니다.
// 모든 Lambda 함수의 응답이 이 형식을 따르도록 하여 일관성을 유지합니다.
type apiResponse struct {
	StatusCode int               `json:"statusCode"` // HTTP 상태 코드 (예: 200, 201, 400, 500)
	Headers    map[string]string `json:"headers"`    // HTTP 응답 헤더
	Body       string            `json:"body"`       // 응답 본문 (일반적으로 JSON 문자열)
}

// apiGatewayResponse 함수는 apiResponse 구조체를 쉽게 생성하기 위한 헬퍼 함수입니다.
// 공통적인 HTTP 헤더(특히 CORS 관련 헤더)를 설정하여 클라이언트 측에서 발생할 수 있는 문제를 방지합니다.
func apiGatewayResponse(statusCode int, body string) apiResponse {
	return apiResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                 "application/json",         // 응답 본문이 JSON임을 명시
			"Access-Control-Allow-Origin":  "*",                        // 모든 도메인에서의 요청을 허용 (보안 강화를 위해 실제 환경에서는 특정 도메인으로 제한 권장)
			"Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,OPTIONS", // 허용할 HTTP 메서드
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token", // 허용할 요청 헤더
		},
		Body: body,
	}
}

// handleRequest 함수는 createPost Lambda 함수의 메인 핸들러입니다.
// AWS Lambda는 이 함수를 호출하여 API Gateway로부터 전달된 이벤트를 처리합니다.
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (apiResponse, error) {
	// CloudWatch Logs에 함수 호출과 요청 본문을 기록하여 디버깅에 활용합니다.
	fmt.Printf("createPost 함수 호출됨. 요청 본문: %s\n", request.Body)

	// 1. HTTP 요청 본문 파싱 (JSON 언마샬링)
	// 클라이언트로부터 전송된 JSON 요청 본문을 Go의 임시 구조체에 매핑합니다.
	// ⭐ 이 newPostData 구조체에 Author 필드를 추가해야 합니다. ⭐
	var newPostData struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Author  string `json:"author"` // ⭐ 이 줄을 추가합니다. ⭐
	}
	err := json.Unmarshal([]byte(request.Body), &newPostData)
	if err != nil {
		// JSON 파싱 오류 발생 시, 400 Bad Request 응답을 반환합니다.
		fmt.Printf("요청 본문 파싱 오류: %v\n", err)
		return apiGatewayResponse(400, `{"message": "Invalid request body"}`), nil
	}

	// 입력 데이터 유효성 검사: 제목, 내용, 작성자가 비어있는지 확인합니다.
	// ⭐ Author 유효성 검사 조건을 추가해야 합니다. ⭐
	if newPostData.Title == "" || newPostData.Content == "" || newPostData.Author == "" {
		fmt.Println("제목, 내용 또는 작성자가 비어 있습니다.")
		return apiGatewayResponse(400, `{"message": "Title, content, and author cannot be empty"}`), nil
	}

	// 2. 새로운 Post 객체 생성 및 필드 초기화
	// time.Now().UTC().Format(time.RFC3339)는 현재 UTC 시각을 ISO 8601 형식의 문자열로 변환합니다.
	now := time.Now().UTC().Format(time.RFC3339)
	post := Post{
		PostID:    uuid.New().String(), // `github.com/google/uuid` 패키지를 사용하여 고유한 UUID를 생성합니다.
		Title:     newPostData.Title,
		Content:   newPostData.Content,
		Author:    newPostData.Author, // ⭐ newPostData에서 Author 값을 할당합니다. ⭐
		CreatedAt: now,                // 생성 시각 설정
		UpdatedAt: now,                // 초기 업데이트 시각은 생성 시각과 동일하게 설정
	}

	// 3. Post 객체를 DynamoDB Item 형식으로 마샬링
	// `dynamodbattribute.MarshalMap` 함수는 Go 구조체를 DynamoDB가 이해할 수 있는 속성 값(AttributeValue) 맵으로 변환합니다.
	item, err := dynamodbattribute.MarshalMap(post)
	if err != nil {
		// 마샬링 오류 발생 시, 500 Internal Server Error 응답을 반환합니다.
		fmt.Printf("DynamoDB 항목 마샬링 오류: %v\n", err)
		return apiGatewayResponse(500, `{"message": "Failed to prepare item for storage"}`), nil
	}

	// 4. DynamoDB PutItem API 호출
	// `PutItemInput`은 DynamoDB에 항목을 저장하기 위한 요청 파라미터를 정의합니다.
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName), // `aws.String`은 Go 문자열을 AWS SDK의 *string 타입으로 변환합니다.
		Item:      item,                  // 마샬링된 DynamoDB 항목 데이터
	}

	// DynamoDB 클라이언트를 사용하여 PutItem 작업을 수행합니다.
	_, err = dynamoClient.PutItem(input)
	if err != nil {
		// DynamoDB 저장 중 오류 발생 시, 500 Internal Server Error 응답을 반환합니다.
		fmt.Printf("DynamoDB PutItem 오류: %v\n", err)
		return apiGatewayResponse(500, `{"message": "Failed to create post"}`), nil
	}

	fmt.Printf("게시물 성공적으로 생성됨: PostID=%s\n", post.PostID)

	// 5. 성공 응답 반환
	// 생성된 Post 객체를 JSON 문자열로 다시 마샬링하여 클라이언트에게 응답 본문으로 전송합니다.
	responseBody, err := json.Marshal(post)
	if err != nil {
		fmt.Printf("응답 본문 마샬링 오류: %v\n", err)
		return apiGatewayResponse(500, `{"message": "Failed to marshal response"}`), nil
	}

	// 201 Created HTTP 상태 코드와 함께 성공 응답을 반환합니다.
	return apiGatewayResponse(201, string(responseBody)), nil
}

func main() {
	// Lambda 핸들러 함수 시작
	lambda.Start(handleRequest)
}
