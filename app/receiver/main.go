package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/vitalvas/owntracks-serverless/internal/message"
)

type Handler struct {
	sqs *sqs.Client

	locationQueue string
}

type Message struct {
	Type string `json:"_type"`
}

func (h *Handler) Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var data Message

	body := request.Body

	if request.IsBase64Encoded {
		decodedBody, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       err.Error(),
			}, nil
		}

		body = string(decodedBody)
	}

	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	switch data.Type {
	case "location":
		if err := h.handleLocation(ctx, request, body); err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       err.Error(),
			}, nil
		}

	default:
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusNotImplemented,
			Body:       http.StatusText(http.StatusNotImplemented),
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "[]",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func (h *Handler) handleLocation(ctx context.Context, request events.APIGatewayV2HTTPRequest, body string) error {
	var data message.MessageLocation

	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return err
	}

	if row, ok := request.Headers["x-limit-u"]; ok {
		data.UserID = row
	}

	if row, ok := request.Headers["x-limit-d"]; ok {
		data.DeviceID = row
	}

	if row, ok := request.Headers["x-forwarded-for"]; ok {
		data.RemoteIP = getRemoteIP(row)
	}

	if len(data.UserID) <= 2 || len(data.DeviceID) <= 2 {
		var unstruct map[string]string

		if err := json.Unmarshal([]byte(body), &unstruct); err != nil {
			return err
		}

		if topic, ok := unstruct["topic"]; ok {
			if row := strings.Split(topic, "/"); len(row) == 3 {
				if len(data.UserID) <= 2 {
					data.UserID = row[1]
				}

				if len(data.DeviceID) <= 2 {
					data.DeviceID = row[2]
				}
			}
		}
	}

	dataByte, err := json.Marshal(data)
	if err != nil {
		return err
	}

	sqsMessage := &sqs.SendMessageInput{
		QueueUrl:    aws.String(h.locationQueue),
		MessageBody: aws.String(string(dataByte)),
	}

	if _, err := h.sqs.SendMessage(ctx, sqsMessage); err != nil {
		return err
	}

	return nil
}

func main() {
	h := Handler{
		locationQueue: os.Getenv("SQS_LOCATION_URL"),
	}

	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	h.sqs = sqs.NewFromConfig(sdkConfig)

	lambda.Start(h.Handler)
}

func getRemoteIP(data string) string {
	for _, row := range strings.Split(data, ",") {
		if netIP := net.ParseIP(row); netIP != nil {
			return row
		}
	}

	return ""
}
