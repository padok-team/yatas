package lambda

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func GetLambdas(s aws.Config) []types.FunctionConfiguration {
	svc := lambda.NewFromConfig(s)
	input := &lambda.ListFunctionsInput{
		MaxItems: aws.Int32(100),
	}
	result, err := svc.ListFunctions(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	return result.Functions
}

func CheckIfLambdaPrivate(s aws.Config, lambdas []types.FunctionConfiguration, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "Lambda Private"
	check.Id = testName
	check.Description = "Check if all Lambdas are private"
	check.Status = "OK"
	for _, lambda := range lambdas {
		if lambda.VpcConfig == nil {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Lambda " + *lambda.FunctionName + " is public"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *lambda.FunctionArn})
		} else {
			status := "OK"
			Message := "Lambda " + *lambda.FunctionName + " is private"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *lambda.FunctionArn})
		}
	}
	*c = append(*c, check)
}

func CheckIfLambdaInSecurityGroup(s aws.Config, lambdas []types.FunctionConfiguration, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "Lambda In Security Group"
	check.Id = testName
	check.Description = "Check if all Lambdas are in a security group"
	check.Status = "OK"
	for _, lambda := range lambdas {
		if lambda.VpcConfig == nil || lambda.VpcConfig.SecurityGroupIds == nil {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Lambda " + *lambda.FunctionName + " is not in a security group"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *lambda.FunctionArn})
		} else {
			status := "OK"
			Message := "Lambda " + *lambda.FunctionName + " is in a security group"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *lambda.FunctionArn})
		}
	}
	*c = append(*c, check)
}

func RunLambdaTests(s aws.Config, c *yatas.Config) []results.Check {
	var checks []results.Check
	lambdas := GetLambdas(s)
	yatas.CheckTest(c, "AWS_LMD_001", CheckIfLambdaPrivate)(s, lambdas, "AWS_LMD_001", &checks)
	yatas.CheckTest(c, "AWS_LMD_002", CheckIfLambdaInSecurityGroup)(s, lambdas, "AWS_LMD_002", &checks)
	return checks
}
