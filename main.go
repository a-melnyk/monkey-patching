package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
)

type IMDSClient interface {
	GetInstanceIdentityDocument(ctx context.Context, params *imds.GetInstanceIdentityDocumentInput, optFns ...func(*imds.Options)) (*imds.GetInstanceIdentityDocumentOutput, error)
}

var NewIMDSClient func(cfg aws.Config) IMDSClient = func(cfg aws.Config) IMDSClient {
	return imds.NewFromConfig(cfg)
}

func CallAws(client http.Client) (string, error) {
	awsConfig, err := awscfg.LoadDefaultConfig(context.TODO(), awscfg.WithHTTPClient(&client))
	if err != nil {
		return "", err
	}

	awsClient := NewIMDSClient(awsConfig)
	instanceIdentity, err := awsClient.GetInstanceIdentityDocument(context.TODO(), &imds.GetInstanceIdentityDocumentInput{})
	if err != nil {
		return "", err
	}

	return instanceIdentity.InstanceIdentityDocument.InstanceID, nil
}

func main() {
	println("Starting main")
	println(CallAws(http.Client{}))
}
