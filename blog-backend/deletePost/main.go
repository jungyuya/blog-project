// deletePost/main.go
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

	"blog-backend/utils"
)

var (
	sess         = session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	dynamoClient = dynamodb.New(sess)
	tableName    = os.Getenv("DYNAMODB_TABLE")
)

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	origin := utils.ValidateOrigin(req.Headers["origin"])

	postID := req.PathParameters["postId"]
	if len(postID) == 0 {
		return utils.NewCORSResponse(400, `{"message":"PostID is required"}`, origin), nil
	}

	_, err := dynamoClient.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       map[string]*dynamodb.AttributeValue{"postId": {S: aws.String(postID)}},
	})
	if err != nil {
		fmt.Printf("DynamoDB DeleteItem error: %v\n", err)
		return utils.NewCORSResponse(500, `{"message":"Failed to delete post"}`, origin), nil
	}

	fmt.Printf("PostID %s deleted successfully\n", postID)
	return utils.NewCORSResponse(204, "", origin), nil
}

func main() {
	lambda.Start(handleRequest)
}
