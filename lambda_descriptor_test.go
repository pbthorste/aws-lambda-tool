package lambda_deploy
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"io/ioutil"
)

func TestLoadDescriptor(t *testing.T) {
	fileName := "./testdata/test1/lambda-desc.yml"
	lambdaDesc := LoadDescriptorFile(fileName)
	assert.Equal(t, lambdaDesc.Function_name, "python-hello", "should be equal")
}

func TestLoadDescriptorNoVpc(t *testing.T) {
	fileName := "./testdata/test1/lambda-desc.yml"
	lambdaDesc := LoadDescriptorFile(fileName)
	assert.Nil(t,lambdaDesc.Vpc_config, "should be nil")
}

func TestLoadDescriptorWithVpc(t *testing.T) {
	fileName := "./testdata/descriptors/vpc-descriptor.yml"
	lambdaDesc := LoadDescriptorFile(fileName)
	assert.NotNil(t,lambdaDesc.Vpc_config, "should not be nil")
	assert.Len(t, (lambdaDesc.Vpc_config).Subnet_ids, 2)
	assert.Len(t, (lambdaDesc.Vpc_config).Security_group_ids, 2)
}

func TestLoadDescriptorWithBadVpc1(t *testing.T) {
	fileName := "./testdata/descriptors/vpc-descriptor-bad1.yml"
	data, err := ioutil.ReadFile(fileName)
	check(err)
	lambdaParent := unmarshalDescriptor(data)
	error := lambdaParent.Lambda.Validate()
	assert.Error(t, error, "There should be an error ")
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


func TestCompareVpcDescriptorHasNone(t *testing.T) {
	lambdaDesc := LambdaFunctionDesc{Function_name:"my-function"}
	response := lambda.VpcConfigResponse{}
	_, isDifferent := lambdaDesc.CompareVpcConfig(&response)
	assert.True(t, isDifferent, "Should be different")
}


func TestCompareVpcBothHave(t *testing.T) {
	lambdaDesc := LambdaFunctionDesc{Function_name:"my-function"}
	newVpc := LambdaVpcConfig{}
	lambdaDesc.Vpc_config = &newVpc
	response := lambda.VpcConfigResponse{}
	_, isDifferent := lambdaDesc.CompareVpcConfig(&response)
	assert.False(t, isDifferent, "Should not be different")
}

func TestCompareVpcBothHaveSameValue(t *testing.T) {
	lambdaDesc := LambdaFunctionDesc{Function_name:"my-function"}
	newVpc := LambdaVpcConfig{
		Security_group_ids: []string{"sg1"},
		Subnet_ids: []string{"sub1"},
	}
	lambdaDesc.Vpc_config = &newVpc
	response := lambda.VpcConfigResponse{
		SecurityGroupIds: aws.StringSlice([]string{"sg1"}),
		SubnetIds: aws.StringSlice([]string{"sub1"}),
	}

	_, isDifferent := lambdaDesc.CompareVpcConfig(&response)
	assert.False(t, isDifferent, "Should not be different")
}

func TestCompareVpcBothHaveDiffSg(t *testing.T) {
	lambdaDesc := LambdaFunctionDesc{Function_name:"my-function"}
	newVpc := LambdaVpcConfig{
		Security_group_ids: []string{"sg1"},
		Subnet_ids: []string{"sub1"},
	}
	lambdaDesc.Vpc_config = &newVpc
	response := lambda.VpcConfigResponse{
		SecurityGroupIds: aws.StringSlice([]string{"sg2"}),
		SubnetIds: aws.StringSlice([]string{"sub1"}),
	}

	_, isDifferent := lambdaDesc.CompareVpcConfig(&response)
	assert.True(t, isDifferent, "Should be different")
}

func TestCompareVpcBothHaveDiffSub(t *testing.T) {
	lambdaDesc := LambdaFunctionDesc{Function_name:"my-function"}
	newVpc := LambdaVpcConfig{
		Security_group_ids: []string{"sg1"},
		Subnet_ids: []string{"sub1"},
	}
	lambdaDesc.Vpc_config = &newVpc
	response := lambda.VpcConfigResponse{
		SecurityGroupIds: aws.StringSlice([]string{"sg1"}),
		SubnetIds: aws.StringSlice([]string{"sub2"}),
	}

	_, isDifferent := lambdaDesc.CompareVpcConfig(&response)
	assert.True(t, isDifferent, "Should be different")
}
