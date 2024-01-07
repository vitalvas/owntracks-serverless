package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Handler struct {
	secretName string
	sec        *secretsmanager.Client
}

func (h *Handler) HandleRequest(ctx context.Context, request events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	token, exists := request.Headers["authorization"]
	if !exists {
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	// extract userpass from "basic <userpass base64>" string
	tokenSlice := strings.Split(token, " ")
	if len(tokenSlice) != 2 {
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	var (
		login string
		pass  string
	)

	switch strings.ToLower(tokenSlice[0]) {
	case "basic":
		decoded, err := base64.StdEncoding.DecodeString(tokenSlice[1])
		if err != nil {
			return events.APIGatewayV2CustomAuthorizerSimpleResponse{
				IsAuthorized: false,
			}, nil
		}

		loginPass := strings.Split(string(decoded), ":")

		if len(loginPass) != 2 {
			return events.APIGatewayV2CustomAuthorizerSimpleResponse{
				IsAuthorized: false,
			}, nil
		}

		login = loginPass[0]
		pass = loginPass[1]

	default:
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(h.secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := h.sec.GetSecretValue(ctx, input)
	if err != nil {
		log.Println(err)

		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	var authData map[string]string
	if err = json.Unmarshal([]byte(*result.SecretString), &authData); err != nil {
		log.Println(err)

		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	if authPass, ok := authData[login]; ok {
		if authPass == pass {
			return events.APIGatewayV2CustomAuthorizerSimpleResponse{
				IsAuthorized: true,
				Context: map[string]interface{}{
					"username": login,
				},
			}, nil
		}
	}

	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: false,
	}, nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	h := Handler{
		secretName: os.Getenv("SECRET_ID"),
		sec:        secretsmanager.NewFromConfig(sdkConfig),
	}

	lambda.Start(h.HandleRequest)
}
