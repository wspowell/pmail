package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/wspowell/snailmail/resources"
	"github.com/wspowell/snailmail/resources/db"
)

const (
	headerOneTimePassword = "X-OneTimePassword"
)

type connect struct {
	Datastore db.Datastore
}

func (self *connect) handler(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx := context.Background()

	// Because we're using cloudwatch for logs, printlines work as logs
	fmt.Println("Begin Connection")
	fmt.Println("Connection ID: " + request.RequestContext.ConnectionID)

	oneTimePassword, exists := request.Headers[headerOneTimePassword]
	if !exists {
		fmt.Printf("missing header '%s'", headerOneTimePassword)
		return events.APIGatewayProxyResponse{
			StatusCode: 422,
		}, nil
	}

	if err := self.Datastore.CreateOneTimePassword(ctx, request.RequestContext.ConnectionID, oneTimePassword); err != nil {
		fmt.Printf("failed to create OTP: %+v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, nil
	}

	// Have to return this for the lambda to understand what happened. If you don't, the lambda will be confused and send odd errors to the caller
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	apiResources := resources.Load()

	handle := &connect{
		Datastore: apiResources.Datastore,
	}

	lambda.Start(handle.handler)
}
