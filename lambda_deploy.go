package lambda_deploy

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/mitchellh/go-homedir"
	"fmt"
	"strings"
	"io/ioutil"
	"log"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"crypto/sha256"
	"encoding/base64"
)



func LambdaDeploy(profile, region, zipfile string, descriptor *LambdaFunctionDesc) {
	svc := SetupLambdaClient(profile, region)
	getFunctionInput := lambda.GetFunctionInput{FunctionName: &(descriptor.Function_name)}
	result, err := svc.GetFunction(&getFunctionInput)

	if checkIfLambdaIsDeployed(err) {
		fmt.Println("The function already exists")
		amazonSha := *(result.Configuration.CodeSha256)
		if amazonSha == Base64sha256(zipfile) {
			fmt.Println("Your zipfile and the uploaded one are identical")
		} else {
			fmt.Println("Uploading lambda function")
			updateExistingCode(svc, descriptor, zipfile)
		}
		configDiff, isDifferent := descriptor.CompareConfig(result.Configuration)
		if !isDifferent {
			fmt.Println("Config is unchanged - will not update")
		} else {
			fmt.Println("Config is changed - differences:", configDiff)
			result, err := svc.UpdateFunctionConfiguration(configDiff)
			check(err)
			fmt.Println("Config has been updated, it is now:", result)
		}
	} else {
		fmt.Println("Lambda function is not deployed")
		createNewLambda(svc, descriptor, zipfile)
	}
}

func checkIfLambdaIsDeployed(getFunctionError error) (bool) {
	if getFunctionError == nil {
		return true
	} else {
		if strings.Contains(getFunctionError.Error(), "ResourceNotFoundException") {
			return false
		} else {
			panic(getFunctionError)
		}
	}
}

func createNewLambda(client *lambda.Lambda, descriptor *LambdaFunctionDesc, zipfile string) {
	file, zipErr := loadFileContent(zipfile)
	if zipErr != nil {
		panic(fmt.Errorf("Unable to load %q: %s", zipfile, zipErr))
	}
	functionCode := lambda.FunctionCode {
		ZipFile: file,
	}
	params := &lambda.CreateFunctionInput{
		Code:         &functionCode,
		Description:  aws.String(descriptor.Description),
		FunctionName: aws.String(descriptor.Function_name),
		Handler:      aws.String(descriptor.Handler),
		MemorySize:   aws.Int64(int64(descriptor.Memory_size)),
		Role:         aws.String(descriptor.Role),
		Runtime:      aws.String(descriptor.Runtime),
		Timeout:      aws.Int64(int64(descriptor.Timeout)),
		Publish:      aws.Bool(descriptor.Publish),
	}
	if len(descriptor.Environment) > 0 {
		params.Environment = &lambda.Environment{
			Variables: aws.StringMap(descriptor.Environment),
		}
	}
	if descriptor.Vpc_config != nil {
		params.VpcConfig = &lambda.VpcConfig{
			SecurityGroupIds: aws.StringSlice(descriptor.Vpc_config.Security_group_ids),
			SubnetIds: aws.StringSlice(descriptor.Vpc_config.Subnet_ids),
		}
	}
	fmt.Println("Uploading lambda function")
	_, err := client.CreateFunction(params)
	if err != nil {
		log.Printf("[ERROR] Received %q", err)
		if awserr, ok := zipErr.(awserr.Error); ok {
			if awserr.Code() == "InvalidParameterValueException" {
				log.Printf("[DEBUG] InvalidParameterValueException creating Lambda Function: %s", awserr)

			}
		}
		log.Printf("[DEBUG] Error creating Lambda Function: %s", err)
		panic(err)
	}
}

func updateExistingCode(client *lambda.Lambda, descriptor *LambdaFunctionDesc, zipfile string) {
	file, err := loadFileContent(zipfile)
	check(err)
	input := &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(descriptor.Function_name),
		Publish:      aws.Bool(descriptor.Publish),
		ZipFile:      file,
	}
	result, err := client.UpdateFunctionCode(input)
	check(err)
	fmt.Println("Result of code update:", result)
}

// see: https://github.com/hashicorp/terraform/blob/master/builtin/providers/aws/resource_aws_lambda_function.go
func loadFileContent(v string) ([]byte, error) {
	filename, err := homedir.Expand(v)
	if err != nil {
		return nil, err
	}
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}


// Generates a base64 sha 256 string from a file in order to check if it is
// Different from the one that is on AWS.
func Base64sha256 (zipfile string) (string) {
	file, _ := loadFileContent(zipfile)
	h := sha256.New()
	h.Write(file)
	shaSum := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(shaSum[:])
}

