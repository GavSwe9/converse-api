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
	"encoding/json"
)

type MessageItem struct {
	GroupId  int    `json:"groupId"`
	UserName string `json:"userName"`
	Message  string `json:"message"`
}

func (e *MessageItem) Decode(data []byte) (*MessageItem, error) {
	err := json.Unmarshal(data, e)
	return e, err
}

type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (Response, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	fmt.Println(ctx)
	fmt.Println(request.Body)

	svc := dynamodb.New(sess)
	
	messageItem, err := new(MessageItem).Decode([]byte(request.Body))

	dynamoItem, err := dynamodbattribute.MarshalMap(messageItem)
	if err != nil {
		log.Fatalf("Error marshalling new message item: %s", err)
	}

	tableName := os.Getenv("MESSAGES_TABLE")

	input := &dynamodb.PutItemInput{
		Item:      dynamoItem,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Error calling PutItem: %s", err)
	}

	fmt.Println("Successfully put message")

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
