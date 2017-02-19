package lambda_deploy

import (
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
)

func SetupLambdaClient(profile, region string) (*lambda.Lambda) {
	config := aws.NewConfig()
	if region != "" {
		config = config.WithRegion(region)
	}
	options := session.Options{
		Config:            *config,
		SharedConfigState: session.SharedConfigEnable,
	}
	if profile != "" {
		options.Profile = profile
	}
	sess, err := session.NewSessionWithOptions(options)
	check(err)
	return lambda.New(sess)
}
