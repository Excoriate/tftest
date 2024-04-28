package cloudprovider

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type AWSAdapter interface {
	NewSNS() *sns.Client
	NewSQS() *sqs.Client
	NewS3() *s3.Client
	NewRDS() *rds.Client
	NewEC2() *ec2.Client
	NewIAM() *iam.Client
	NewDynamoDB() *dynamodb.Client
	NewAutoScaling() *autoscaling.Client
	NewECS() *ecs.Client
	NewEKS() *eks.Client
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

func (a *AWS) NewS3() *s3.Client {
	return s3.NewFromConfig(a.cfg)
}

func (a *AWS) NewRDS() *rds.Client {
	return rds.NewFromConfig(a.cfg)
}

func (a *AWS) NewEC2() *ec2.Client {
	return ec2.NewFromConfig(a.cfg)
}

func (a *AWS) NewIAM() *iam.Client {
	return iam.NewFromConfig(a.cfg)
}

func (a *AWS) NewDynamoDB() *dynamodb.Client {
	return dynamodb.NewFromConfig(a.cfg)
}

func (a *AWS) NewAutoScaling() *autoscaling.Client {
	return autoscaling.NewFromConfig(a.cfg)
}

func (a *AWS) NewECS() *ecs.Client {
	return ecs.NewFromConfig(a.cfg)
}

func (a *AWS) NewEKS() *eks.Client {
	return eks.NewFromConfig(a.cfg)
}
