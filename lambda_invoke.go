package lambda_deploy

import (
	"github.com/aws/aws-sdk-go/service/lambda"
)

func InvokeLambda(client *lambda.Lambda, functionName, body string) (string) {
	invoke := lambda.InvokeInput{}
	invoke.SetFunctionName(functionName)
	if body != "" {
		invoke.SetPayload([]byte(body))
	}
	out, err := client.Invoke(&invoke)
	check(err)
	//fmt.Println(out)
	return string(out.Payload)
}
