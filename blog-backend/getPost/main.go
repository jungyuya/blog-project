// getPost/main.go
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

	"blog-backend/utils"
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

	postID := req.PathParameters["postId"]
	if len(postID) == 0 {
		return utils.NewCORSResponse(400, `{"message":"PostID is required"}`, origin), nil
	}

	out, err := dynamoClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       map[string]*dynamodb.AttributeValue{"postId": {S: aws.String(postID)}},
	})
	if err != nil {
		fmt.Printf("DynamoDB GetItem error: %v\n", err)
		return utils.NewCORSResponse(500, `{"message":"Failed to retrieve post"}`, origin), nil
	}
	if len(out.Item) == 0 {
		return utils.NewCORSResponse(404, `{"message":"Post not found"}`, origin), nil
	}

	var post Post
	if err := dynamodbattribute.UnmarshalMap(out.Item, &post); err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
		return utils.NewCORSResponse(500, `{"message":"Failed to process retrieved post"}`, origin), nil
	}

	body, _ := json.Marshal(post)
	return utils.NewCORSResponse(200, string(body), origin), nil
}

func main() {
	lambda.Start(handleRequest)
}
