package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"context"
	"fmt"
	"log"
	"os"
)

type ConnectionItem struct {
	ConnectionId string `json:"connectionId"`
	UserName     string `json:"userName"`
	GroupId      int    `json:"groupId"`
}

type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (Response, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	item := ConnectionItem{
		ConnectionId: request.RequestContext.ConnectionID,
		UserName:     "Gavin.Sweeney",
		GroupId:      1,
	}

	dynamoItem, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Error marshalling new connection item: %s", err)
	}

	tableName := os.Getenv("CONNECTIONS_TABLE")

	input := &dynamodb.PutItemInput{
		Item:      dynamoItem,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Error calling PutItem: %s", err)
	}

	fmt.Println("Successfully added connection to table")

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	fmt.Println(resp)
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
