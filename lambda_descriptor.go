package lambda_deploy
/*
Processes descriptor for lambda
 */

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"github.com/aws/aws-sdk-go/service/lambda"
	"errors"
	"strings"
	"github.com/aws/aws-sdk-go/aws"
)

func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
		panic(e)
	}
}

type LambdaDescriptor struct {
	Lambda LambdaFunctionDesc
}

type LambdaFunctionDesc struct {
	Function_name string
	Description string
	Handler string
	Runtime string
	Role string
	Memory_size int
	Timeout int
	Publish bool  //default for bool is false, which fits in this case
	Environment map[string]string
	Vpc_config *LambdaVpcConfig
}

type LambdaVpcConfig struct {
	Subnet_ids []string
	Security_group_ids []string
}

func (l *LambdaFunctionDesc) SetDefaults() {
	if l.Timeout == 0 {
		l.Timeout = 3
	}
	if l.Memory_size == 0 {
		l.Memory_size = 128
	}
}

func (l *LambdaFunctionDesc) Validate() error {
	errorList := make([]string, 0)
	if l.Function_name == "" {
		errorList = append(errorList, "Missing function_name")
	}
	if l.Handler == "" {
		errorList = append(errorList, "Missing handler")
	}
	if l.Runtime == "" {
		errorList = append(errorList, "Missing runtime")
	}
	if l.Role == "" {
		errorList = append(errorList, "Missing role")
	}
	if l.Vpc_config != nil {
		if len(l.Vpc_config.Security_group_ids) < 1 {
			errorList = append(errorList, "There must be at least 1 vpc security group id")
		}
		if len(l.Vpc_config.Subnet_ids) < 1 {
			errorList = append(errorList, "There must be at least 1 vpc subnet id")
		}
	}
	if len(errorList) > 0 {
		return errors.New("Descriptor error: " + strings.Join(errorList, ","))
	}
	return nil
}

func (d *LambdaFunctionDesc) CompareConfig(functionConfig *lambda.FunctionConfiguration) (*lambda.UpdateFunctionConfigurationInput, bool) {
	isDifferent := false
	input := lambda.UpdateFunctionConfigurationInput{}
	input.SetFunctionName(d.Function_name)
	if (functionConfig.Description == nil && d.Description != "") ||
		(functionConfig.Description != nil && *functionConfig.Description != d.Description){
		input.SetDescription(d.Description)
		isDifferent = true
	}
	if (functionConfig.Handler == nil && d.Handler != "") ||
		(functionConfig.Handler != nil && *functionConfig.Handler != d.Handler) {
		input.SetHandler(d.Handler)
		isDifferent = true
	}
	if (functionConfig.MemorySize == nil && d.Memory_size != 0) ||
		(functionConfig.MemorySize != nil && *functionConfig.MemorySize != int64(d.Memory_size)) {
		input.SetMemorySize(int64(d.Memory_size))
		isDifferent = true
	}
	if (functionConfig.Timeout == nil && d.Timeout != 0) ||
	        (functionConfig.Timeout != nil && *functionConfig.Timeout != int64(d.Timeout)) {
		input.SetTimeout(int64(d.Timeout))
		isDifferent = true
	}
	if (functionConfig.Role == nil && d.Role != "") ||
		(functionConfig.Role != nil && *functionConfig.Role != d.Role) {
		input.SetRole(d.Role)
		isDifferent = true
	}
	if (functionConfig.Runtime == nil && d.Runtime != "") ||
		(functionConfig.Runtime != nil && *functionConfig.Runtime != d.Runtime) {
		input.SetRuntime(d.Runtime)
		isDifferent = true
	}
	if (functionConfig.Environment == nil && len(d.Environment) != 0) {
		input.SetEnvironment(&lambda.Environment{
			Variables: aws.StringMap(d.Environment),
		})
		isDifferent = true
	}
	if(functionConfig.Environment != nil){
		if newEnv, isDiff := d.CompareEnvironmentConfig(functionConfig.Environment); isDiff {
			input.SetEnvironment(&lambda.Environment{
				Variables: newEnv,
			})
			isDifferent = true
		}
	}
	if(functionConfig.VpcConfig == nil && d.Vpc_config != nil) {
		input.SetVpcConfig(&lambda.VpcConfig{
			SecurityGroupIds: aws.StringSlice(d.Vpc_config.Security_group_ids),
			SubnetIds: aws.StringSlice(d.Vpc_config.Subnet_ids),
		})
		isDifferent = true
	}
	if(functionConfig.VpcConfig != nil) {
		if newVpc, isDiff := d.CompareVpcConfig(functionConfig.VpcConfig); isDiff {
			input.SetVpcConfig(newVpc)
			isDifferent = true
		}
	}
	err := input.Validate()
	check(err)
	return &input, isDifferent
}

func (d *LambdaFunctionDesc) CompareEnvironmentConfig(other *lambda.EnvironmentResponse) (map[string]*string, bool) {
	isDifferent := false
	if len(d.Environment) != len(other.Variables) {
		isDifferent = true
	} else {
		for k, v := range d.Environment {
			if val, ok := other.Variables[k]; ok {
				if v != *val {
					isDifferent = true
					break
				}
			} else {
				isDifferent = true
				break
			}
		}
	}
	if !isDifferent {
		return nil, isDifferent
	} else {
		return aws.StringMap(d.Environment), isDifferent
	}
}

func (d *LambdaFunctionDesc) CompareVpcConfig(other *lambda.VpcConfigResponse) (*lambda.VpcConfig, bool) {
	isDifferent := false
	if other == nil && d.Vpc_config == nil {
		isDifferent = false
	} else if other != nil && d.Vpc_config == nil {
		// returning empty config will erase it on aws
		return &lambda.VpcConfig{}, true
	} else if len(d.Vpc_config.Subnet_ids) != len(other.SubnetIds) ||
	   len(d.Vpc_config.Security_group_ids) != len(other.SecurityGroupIds) {
		isDifferent = true
	}
	if len(d.Vpc_config.Subnet_ids) == len(other.SubnetIds) &&
		len(d.Vpc_config.Security_group_ids) == len(other.SecurityGroupIds) {

		if !compareSlices(d.Vpc_config.Subnet_ids, other.SubnetIds ) {
			isDifferent = true
		}
		if !compareSlices(d.Vpc_config.Security_group_ids, other.SecurityGroupIds ) {
			isDifferent = true
		}
	}

	if !isDifferent {
		return nil, isDifferent
	} else {
		return &lambda.VpcConfig{
			SecurityGroupIds: aws.StringSlice(d.Vpc_config.Security_group_ids),
			SubnetIds: aws.StringSlice(d.Vpc_config.Subnet_ids),
		}, isDifferent
	}
}

func compareSlices(descriptor []string, other []*string) (bool) {
	isSame := true
	for _,val1 := range descriptor {
		found := false
		for _, val2 := range other {
			if val1 == *val2 {
				found = true
				break
			}
		}
		if found == false {
			isSame = false
		}
	}
	return  isSame
}

func LoadDescriptorFile(filename string) (*LambdaFunctionDesc) {
	data, err := ioutil.ReadFile(filename)
	check(err)
	return LoadDescriptor(data)
}

func unmarshalDescriptor(contents []byte) (*LambdaDescriptor) {
	lambdaParent := LambdaDescriptor{}
	err := yaml.Unmarshal([]byte(contents), &lambdaParent)
	check(err)
	return &lambdaParent
}

func LoadDescriptor(contents []byte) (*LambdaFunctionDesc) {
	lambdaParent := unmarshalDescriptor(contents)
	lambdaParent.Lambda.SetDefaults()
	err := lambdaParent.Lambda.Validate()
	check(err)
	return &lambdaParent.Lambda
}
