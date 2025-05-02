package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

type CognitoClient struct {
	Region      string
	UserPoolID  string
	AppClientID string
}

func NewCognitoClient(region, userPoolID, appClientID string) (*CognitoClient, error) {
	ctx := context.Background()
	awsCfg, err := config.LoadDefaultConfig(ctx,
		otelaws.WithTracedInputShapes(),
		otelaws.WithClientOptions(),
	)
	if err != nil {
		return nil, err
	}

	return &CognitoClient{
		Region:      region,
		UserPoolID:  userPoolID,
		AppClientID: appClientID,
	}, nil
}
