package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/vitalvas/owntracks-serverless/internal/message"
)

type Handler struct {
	locationTable string
	db            *dynamodb.Client
}

func (h *Handler) Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	if sqsEvent.Records == nil {
		return nil
	}

	items := make([]types.WriteRequest, len(sqsEvent.Records))

	for idx, msg := range sqsEvent.Records {
		var data message.MessageLocation

		if err := json.Unmarshal([]byte(msg.Body), &data); err != nil {
			return err
		}

		av, err := attributevalue.MarshalMap(data)
		if err != nil {
			return err
		}

		items[idx] = types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: av,
			},
		}
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			h.locationTable: items,
		},
	}

	if _, err := h.db.BatchWriteItem(ctx, input); err != nil {
		return err
	}

	return nil
}

func main() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	h := Handler{
		locationTable: os.Getenv("DB_LOCATION_NAME"),
		db:            dynamodb.NewFromConfig(sdkConfig),
	}

	lambda.Start(h.Handler)
}
