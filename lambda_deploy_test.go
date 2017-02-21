package lambda_deploy

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

/*
type LambdaClientMock struct{
	mock.Mock
}
func (m *LambdaClientMock) UpdateFunctionCode(input *lambda.UpdateFunctionCodeInput) (*lambda.FunctionConfiguration, error) {
	args := m.Called(input)
	return *args.Get(0).(lambda.FunctionConfiguration), args.Error(1)
}
*/
func TestBase641(t *testing.T) {
	fileName := "./testdata/test1/python_hello.zip"
	result := Base64sha256(fileName)
	assert.Equal(t, "MqhRu7AvFO9UcpcXI4tzTp63SMLtEm6UQhl54W1w0Ss=", result)
}