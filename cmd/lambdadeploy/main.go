package main

import (
	"gopkg.in/urfave/cli.v1"
	"fmt"
	"os"
	"github.com/pbthorste/aws-lambda-tool"
	"errors"
	"io/ioutil"
)
var (
	version string
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	app := cli.NewApp()
	app.Name = "lambda-deploy (lambdadeploy)"
	app.Version = version
	app.Usage = "Deploys Amazon Lambda Functions"
	app.Flags = []cli.Flag {
		cli.BoolFlag{
			Name: "noheader, H",
			Usage: "Do not show header",
		},
		cli.StringFlag{
			Name: "region",
			Usage: "AWS region (optional, can be set by env var 'AWS_REGION' or in profile)",
		},
		cli.StringFlag{
			Name: "profile",
			Usage: "AWS profile (optional, can be set by env var 'AWS_PROFILE')",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Usage:   "list lambda functions",
			Action:  func(c *cli.Context) error {
				if !c.GlobalBool("noheader") {
					fmt.Println("Installed lambdas\n----------------------")
				}
				client := lambda_deploy.SetupLambdaClient(c.GlobalString("profile"), c.GlobalString("region"))
				lambdas := lambda_deploy.ListLambdas(client)
				fmt.Println(lambdas)
				return nil
			},
		},
		{
			Name:    "delete",
			Usage:   "delete a lambda function",
			Flags:   []cli.Flag{
				cli.StringFlag{
					Name: "name, n",
					Usage: "`Name` of lambda function (required)",

				},
			},
			Action:  func (c *cli.Context) error {
				name, err := checkRequiredArg("name", c.String("name"))
				if err != nil {
					return cli.NewExitError(err, 2)
				}
				if !c.GlobalBool("noheader") {
					fmt.Println("Deleting lambda: " + name + "\n----------------------")
				}
				client := lambda_deploy.SetupLambdaClient(c.GlobalString("profile"), c.GlobalString("region"))
				lambda_deploy.DeleteLambda(client, name)
				fmt.Println("Lambda function has been deleted")
				return nil
			},
		},
		{
			Name: "deploy",
			Usage: "Deploy a lambda function using a descriptor",
			Flags:   []cli.Flag{
				cli.StringFlag{
					Name: "descriptor, d",
					Usage: "`Descriptor` for the lambda function (required)",

				},
				cli.StringFlag{
					Name: "zip-file, z",
					Usage: "`ZIP-File` containing the lambda function (required)",

				},
			},
			Action:  func (c *cli.Context) error {
				descriptor, err := checkRequiredArg("descriptor", c.String("descriptor"))
				if err != nil {
					return cli.NewExitError(err, 2)
				}
				zipfile, err    := checkRequiredArg("zip-file", c.String("zip-file"))
				if err != nil {
					return cli.NewExitError(err, 2)
				}
				lambdaDesc := lambda_deploy.LoadDescriptorFile(descriptor)
				lambda_deploy.LambdaDeploy(c.GlobalString("profile"), c.GlobalString("region"), zipfile, lambdaDesc)
				fmt.Println("Lambda function deployed successfully")
				return nil
			},

		},
		{
			Name: "account",
			Usage: "display account settings",
			Action: func (c *cli.Context) error {
				if !c.GlobalBool("noheader") {
					fmt.Println("Account Settings\n----------------------")
				}
				client := lambda_deploy.SetupLambdaClient(c.GlobalString("profile"), c.GlobalString("region"))
				fmt.Println(lambda_deploy.LambdaAccountSettings(client))
				return nil
			},
		},
		{
			Name: "invoke",
			Usage: "invoke the lambda function synchronously",
			Flags:   []cli.Flag{
				cli.StringFlag{
					Name: "descriptor, d",
					Usage: "`Descriptor` for the lambda function (can not be used with name)",

				},
				cli.StringFlag{
					Name: "name, n",
					Usage: "`Name` of the lambda function (can not be used with descriptor)",

				},
				cli.StringFlag{
					Name: "body, b",
					Usage: "`Body` to send to the lambda function (can not be used with file)",
				},
				cli.StringFlag{
					Name: "file, f",
					Usage: "`File` containing text to be sent to the lambda function (can not be used with body)",
				},
			},
			Action: func (c *cli.Context) error {
				body := c.String("body")
				bodyFile := c.String("file")
				if onlyOne, err := thereMustBeOnlyOne("descriptor", c.String("descriptor"), "name", c.String("name")); !onlyOne {
					return cli.NewExitError(err, 2)
				}
				if onlyOne, err := thereCanBeOnlyOne("body", body, "file", bodyFile); !onlyOne {
					return cli.NewExitError(err, 2)
				}
				functionName := getFunctionName(c.String("descriptor"), c.String("name"))
				if body == "" && bodyFile != "" {
					data, err := ioutil.ReadFile(bodyFile)
					check(err)
					body = string(data)
				}
				if !c.GlobalBool("noheader") {
					fmt.Printf("Invoking lambda function: %v\n----------------------\n", functionName)
				}
				client := lambda_deploy.SetupLambdaClient(c.GlobalString("profile"), c.GlobalString("region"))
				output := lambda_deploy.InvokeLambda(client, functionName, body)
				fmt.Println(output)
				return nil
			},
		},
	}
	app.Run(os.Args)
}

func checkRequiredArg(name, value string) (string, error) {
	if value == "" {
		msg := "Error: missing required argument: " + name
		return "", errors.New(msg)
	} else {
		return value, nil
	}
}

func thereMustBeOnlyOne(firstName, firstValue, secondName, secondValue string) (bool, error) {
	if firstValue == "" && secondValue == "" {
		return false, errors.New(fmt.Sprintf("One of %v or %v must be set", firstName, secondName))
	}
	return thereCanBeOnlyOne(firstName, firstValue, secondName, secondValue)
}
func thereCanBeOnlyOne(firstName, firstValue, secondName, secondValue string) (bool, error) {
	if firstValue != "" && secondValue != "" {
		return false, errors.New(fmt.Sprintf("One of %v or %v can be set - but not both!", firstName, secondName))
	}
	return true, nil
}

// one of descriptor or string must have a value
func getFunctionName(descriptor, name string) (string) {
	if name != "" {
		return name
	} else {
		lambdaDesc := lambda_deploy.LoadDescriptorFile(descriptor)
		return lambdaDesc.Function_name
	}
}