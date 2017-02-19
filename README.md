# aws-lambda-tool
A standalone tool for AWS lambda functions that can among other things deploy and delete them.

# Purpose
This is a CLI tool written in Go that is meant to be able to do simple
operations on AWS lambda functions. 

Currently the available operations are:

* Deploy a lambda
* List deployed lambdas
* Delete a lambda
* Invoke a lambda
* Show lambda statistics


The main purpose of this tool is to deploy lambdas, and be used in a deployment
pipeline.

## Continuous Delivery

This tool can be used to deploy pre-built lambda functions. Provide the tool
the zip/jar file, and a descriptor to deploy the lambda function.



# Installation
The tool can be downloaded from the releases page:

https://github.com/pbthorste/aws-lambda-tool/releases

# Usage
The tool uses amazon profiles as defined by the amazon cli tool (https://aws.amazon.com/cli).
You can either install the cli tool, or manually create the files:

* ~/.aws/config
* ~/.aws/credentials

By default the tool will use the default profile / region.

It can also read this info from environment variables:

* AWS_REGION
* AWS_PROFILE

Or you can give the region/profile as an argument into the script:

```bash
lambdatool --region REGION --profile PROFILE
```

## Deploy a lambda function
You need:
* A zip file containing the lambda function
* A descriptor file

Then you can run:
```bash
lambdatool deploy -d lambda.yml -z lambda.zip
```

# IAM role
Lambda functions need to have an IAM role, and it must be set in the descriptor.
This tool does not create IAM roles - but multiple other tools do, such as:

* [Terraform](https://www.terraform.io/)
* [Ansible](https://www.ansible.com/)

But in order to get started, you could create a basic IAM role that will allow
the lambda function to send logs to cloudwatch. 
Create a role and give it this policy:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:PutLogEvents",
        "logs:CreateLogStream"
      ],
      "Resource": "arn:aws:logs:*:<your account id>:*"
    }
  ]
}
```
You might want to have a more restrictive policy in production environments.

# Descriptor
The tool reads a yaml descriptor with settings for the lambda.
or valid values for the configuration items, see the documentation from Amazon.
http://docs.aws.amazon.com/lambda/latest/dg/welcome.html


filename: lambda.yml
## Java Example
See: 
```json
lambda:
  function_name: aws-lambda-java-example
  description: java hello world
  handler: com.example.lambda.Handler
  runtime: java8
  # TODO: Fix the account number, and create the role
  role: arn:aws:iam::<fix me>:role/basic-lambda-role
```
## Node.js Example
```json
lambda:
  function_name: node-js-app
  description: node-js hello world
  handler: index.handler
  runtime: nodejs4.3
  role: arn:aws:iam::<account id>:role/basic-lambda-role
  memory_size: 128
  timeout: 60
  publish: false
  environment:
    envVar: value
```
## Python Example
```json
lambda:
  function_name: python-hello
  description: python hello world
  handler: python_hello.handler
  runtime: python2.7
  role: arn:aws:iam::<account id>:role/basic-lambda-role
  memory_size: 128
  timeout: 60
  publish: false
  environment:
    envVar: value
```



