package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

    "fmt"
	"os"
    "log"
	"context"
)

type ConnectionKey struct {
    ConnectionId	string	`json:"connectionId"`
}

type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (Response, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	
	svc := dynamodb.New(sess)
	
	key, err := dynamodbattribute.MarshalMap(ConnectionKey{
		request.RequestContext.ConnectionID,
	})
	if err != nil {
		log.Fatalf("Error marshalling connection key: %s", err)
	}

	tableName := os.Getenv("CONNECTIONS_TABLE")

	input := &dynamodb.DeleteItemInput{
		Key:      	key,
		TableName: 	aws.String(tableName),
	}

	_, err = svc.DeleteItem(input)
	if err != nil {
		log.Fatalf("Error calling DeleteItem: %s", err)
	}

	fmt.Println("Successfully removed connection from table")

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type":           "application/json",
		},
	}
	fmt.Println(resp)
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
