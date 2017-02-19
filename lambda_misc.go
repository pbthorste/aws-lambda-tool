package lambda_deploy

import (
	"github.com/aws/aws-sdk-go/service/lambda"
	"strings"
	"fmt"
)

func ListLambdas(client *lambda.Lambda) (string) {
	input := lambda.ListFunctionsInput{}
	resp, err := client.ListFunctions(&input)
	if err != nil {
		if strings.Contains(err.Error(), "NoCredentialProviders") {
			fmt.Println("Error: please check your AWS credentials")
		}
		panic(err)
	}
	return resp.String()
}

func DeleteLambda(client *lambda.Lambda, functionName string) {
	deletionRequest := lambda.DeleteFunctionInput{FunctionName:&functionName}
	_, err := client.DeleteFunction(&deletionRequest)
	check(err)
}


func LambdaAccountSettings(client *lambda.Lambda) (string) {
	input := lambda.GetAccountSettingsInput{}
	result, err := client.GetAccountSettings(&input)
	check(err)
	return result.String()
}
