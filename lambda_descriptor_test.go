package lambda_deploy
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
)

func TestLoadDescriptor(t *testing.T) {
	fileName := "./testdata/test1/lambda-desc.yml"
	lambdaDesc := LoadDescriptorFile(fileName)
	assert.Equal(t, lambdaDesc.Function_name, "python-hello", "should be equal")
}


func TestCompareDescriptor1(t *testing.T) {
	lambdaDesc := LambdaFunctionDesc{Function_name:"my-function"}
	newConfig := lambda.FunctionConfiguration{}
	newConfig.SetFunctionName("my-function")
	_, isDifferent := lambdaDesc.CompareConfig(&newConfig)
	assert.False(t, isDifferent, "should be equal")
}

func TestCompareDescriptor2(t *testing.T) {
	lambdaDesc := LambdaFunctionDesc{Function_name:"my-function"}

	newConfig := lambda.FunctionConfiguration{}
	newConfig.SetFunctionName("my-function")
	envResponse := lambda.EnvironmentResponse{}
	vars := make(map[string]*string)
	vars["yo"] = aws.String("lala")
	envResponse.SetVariables(vars)
	newConfig.SetEnvironment(&envResponse)
	_, isDifferent := lambdaDesc.CompareConfig(&newConfig)
	assert.True(t, isDifferent, "should be different")
}

func TestCompareEnvironmentBothHaveNone(t *testing.T) {
	lambdaDesc := LambdaFunctionDesc{Function_name:"my-function"}
	response := lambda.EnvironmentResponse{}
	_, isDifferent := lambdaDesc.CompareEnvironmentConfig(&response)
	assert.False(t, isDifferent, "Should not be different")
}

func TestCompareEnvironmentOnlyDescriptorHasEnv(t *testing.T) {
	lambdaDesc := LambdaFunctionDesc{Function_name:"my-function"}
	newEnv := make(map[string]string)
	newEnv["key"] = "value"
	lambdaDesc.Environment = newEnv
	response := lambda.EnvironmentResponse{}
	result, isDifferent := lambdaDesc.CompareEnvironmentConfig(&response)
	assert.True(t, isDifferent, "Should be different")
	assert.Equal(t, "value", *result["key"], "should equal")
}

func TestCompareEnvironmentBothHaveSame(t *testing.T) {
	lambdaDesc := LambdaFunctionDesc{Function_name:"my-function"}
	newEnv := make(map[string]string)
	newEnv["key"] = "value"
	lambdaDesc.Environment = newEnv
	response := lambda.EnvironmentResponse{}
	response.Variables = aws.StringMap(newEnv)
	_, isDifferent := lambdaDesc.CompareEnvironmentConfig(&response)
	assert.False(t, isDifferent, "Should not be different")
}

func TestCompareEnvironmentBothHaveSameKeyButDifferentVal(t *testing.T) {
	lambdaDesc := LambdaFunctionDesc{Function_name:"my-function"}
	newEnv := make(map[string]string)
	newEnv["key"] = "value"
	lambdaDesc.Environment = newEnv
	response := lambda.EnvironmentResponse{}
	otherEnv := make(map[string]string)
	otherEnv["key"] = "otherValue"
	response.Variables = aws.StringMap(otherEnv)
	result, isDifferent := lambdaDesc.CompareEnvironmentConfig(&response)
	assert.True(t, isDifferent, "Should be different")
	assert.Equal(t, "value", *result["key"])
}

func TestCompareEnvironmentOtherHasMoreKeys(t *testing.T) {
	lambdaDesc := LambdaFunctionDesc{Function_name:"my-function"}
	newEnv := make(map[string]string)
	newEnv["key"] = "value"
	lambdaDesc.Environment = newEnv
	response := lambda.EnvironmentResponse{}
	otherEnv := make(map[string]string)
	otherEnv["key"] = "value"
	otherEnv["key2"] = "value2"
	response.Variables = aws.StringMap(otherEnv)
	result, isDifferent := lambdaDesc.CompareEnvironmentConfig(&response)
	assert.True(t, isDifferent, "Should be different")
	assert.Equal(t, 1, len(result))
}

func TestCompareEnvironmentOtherOnlyHasKeys(t *testing.T) {
	lambdaDesc := LambdaFunctionDesc{Function_name:"my-function"}
	response := lambda.EnvironmentResponse{}
	otherEnv := make(map[string]string)
	otherEnv["key"] = "value"
	otherEnv["key2"] = "value2"
	response.Variables = aws.StringMap(otherEnv)
	result, isDifferent := lambdaDesc.CompareEnvironmentConfig(&response)
	assert.True(t, isDifferent, "Should be different")
	assert.Equal(t, 0, len(result))
}

func TestCompareValidate1(t *testing.T) {
	empty := LambdaFunctionDesc{}
	err := empty.Validate()
	assert.NotNil(t, err, "There should be an error here")
}
func TestCompareValidate2(t *testing.T) {
	desc := LambdaFunctionDesc{}
	desc.Function_name = "test"
	err := desc.Validate()
	assert.NotNil(t, err, "There should be an error here")
}
