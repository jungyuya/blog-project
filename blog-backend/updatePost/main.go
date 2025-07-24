// updatePost/main.go
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

	"blog-backend/utils"
)

var (
	sess         = session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	dynamoClient = dynamodb.New(sess)
	tableName    = os.Getenv("DYNAMODB_TABLE")
)

type PostUpdate struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	origin := utils.ValidateOrigin(req.Headers["origin"])

	postID := req.PathParameters["postId"]
	if len(postID) == 0 {
		return utils.NewCORSResponse(400, `{"message":"PostID is required"}`, origin), nil
	}

	var input PostUpdate
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		fmt.Printf("Invalid body: %v\n", err)
		return utils.NewCORSResponse(400, `{"message":"Invalid request body"}`, origin), nil
	}
	if input.Title == "" || input.Content == "" || input.Author == "" {
		return utils.NewCORSResponse(400, `{"message":"Title, content, and author cannot be empty"}`, origin), nil
	}

	now := time.Now().UTC().Format(time.RFC3339)
	update := map[string]*dynamodb.AttributeValueUpdate{
		"title": {
			Action: aws.String("PUT"), Value: &dynamodb.AttributeValue{S: aws.String(input.Title)},
		},
		"content": {
			Action: aws.String("PUT"), Value: &dynamodb.AttributeValue{S: aws.String(input.Content)},
		},
		"author": {
			Action: aws.String("PUT"), Value: &dynamodb.AttributeValue{S: aws.String(input.Author)},
		},
		"updatedAt": {
			Action: aws.String("PUT"), Value: &dynamodb.AttributeValue{S: aws.String(now)},
		},
	}

	_, err := dynamoClient.UpdateItem(&dynamodb.UpdateItemInput{
		TableName:        aws.String(tableName),
		Key:              map[string]*dynamodb.AttributeValue{"postId": {S: aws.String(postID)}},
		AttributeUpdates: update,
	})
	if err != nil {
		fmt.Printf("DynamoDB UpdateItem error: %v\n", err)
		return utils.NewCORSResponse(500, `{"message":"Failed to update post"}`, origin), nil
	}

	fmt.Printf("PostID %s updated\n", postID)
	return utils.NewCORSResponse(204, "", origin), nil
}

func main() {
	lambda.Start(handleRequest)
}
