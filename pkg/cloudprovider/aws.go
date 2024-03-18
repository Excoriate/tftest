package cloudprovider

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
)

type AWSAdapter interface {
	NewSNS() *sns.Client
	NewSQS() *sqs.Client
}

type AWS struct {
	Region string
	cfg    aws.Config
}

func NewAWS(region string) (AWSAdapter, error) {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		return nil, fmt.Errorf("AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY must be set")
	}

	cfg, err := awscfg.LoadDefaultConfig(context.TODO(),
		awscfg.WithRegion(region),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}

	return &AWS{Region: region, cfg: cfg}, nil
}

func (a *AWS) NewSNS() *sns.Client {
	return sns.NewFromConfig(a.cfg)
}

func (a *AWS) NewSQS() *sqs.Client {
	return sqs.NewFromConfig(a.cfg)
}
