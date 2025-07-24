// createPost/main.go
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

	"blog-backend/utils"

	"github.com/google/uuid"
)

var (
	sess         = session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	dynamoClient = dynamodb.New(sess)
	tableName    = os.Getenv("DYNAMODB_TABLE")
)

type Post struct {
	PostID    string `json:"postId" dynamodbav:"postId"`
	Title     string `json:"title" dynamodbav:"title"`
	Content   string `json:"content" dynamodbav:"content"`
	Author    string `json:"author" dynamodbav:"author"`
	CreatedAt string `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt string `json:"updatedAt" dynamodbav:"updatedAt"`
}

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	origin := utils.ValidateOrigin(req.Headers["origin"])

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Author  string `json:"author"`
	}
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		fmt.Printf("Invalid body: %v\n", err)
		return utils.NewCORSResponse(400, `{"message":"Invalid request body"}`, origin), nil
	}
	if input.Title == "" || input.Content == "" || input.Author == "" {
		return utils.NewCORSResponse(400, `{"message":"Title, content, and author cannot be empty"}`, origin), nil
	}

	now := time.Now().UTC().Format(time.RFC3339)
	post := Post{
		PostID:    uuid.New().String(),
		Title:     input.Title,
		Content:   input.Content,
		Author:    input.Author,
		CreatedAt: now,
		UpdatedAt: now,
	}

	item, err := dynamodbattribute.MarshalMap(post)
	if err != nil {
		fmt.Printf("Marshal error: %v\n", err)
		return utils.NewCORSResponse(500, `{"message":"Failed to prepare item"}`, origin), nil
	}

	if _, err := dynamoClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}); err != nil {
		fmt.Printf("DynamoDB PutItem error: %v\n", err)
		return utils.NewCORSResponse(500, `{"message":"Failed to create post"}`, origin), nil
	}

	fmt.Printf("Post created: %s\n", post.PostID)
	respBody, _ := json.Marshal(post)
	return utils.NewCORSResponse(201, string(respBody), origin), nil
}

func main() {
	lambda.Start(handleRequest)
}
