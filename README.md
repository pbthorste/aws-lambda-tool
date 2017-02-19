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