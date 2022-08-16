package lambda

import (
	"context"
	"fmt"
	"sync"

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

func CheckIfLambdaPrivate(checkConfig yatas.CheckConfig, lambdas []types.FunctionConfiguration, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Lambda Private", "Check if all Lambdas are private", testName)
	for _, lambda := range lambdas {
		if lambda.VpcConfig == nil {
			Message := "Lambda " + *lambda.FunctionName + " is public"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *lambda.FunctionArn}
			check.AddResult(result)
		} else {
			Message := "Lambda " + *lambda.FunctionName + " is private"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *lambda.FunctionArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

func CheckIfLambdaInSecurityGroup(checkConfig yatas.CheckConfig, lambdas []types.FunctionConfiguration, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Lambda In Security Group", "Check if all Lambdas are in a security group", testName)
	for _, lambda := range lambdas {
		if lambda.VpcConfig == nil || lambda.VpcConfig.SecurityGroupIds == nil {
			Message := "Lambda " + *lambda.FunctionName + " is not in a security group"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *lambda.FunctionArn}
			check.AddResult(result)
		} else {
			Message := "Lambda " + *lambda.FunctionName + " is in a security group"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *lambda.FunctionArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []results.Check
	lambdas := GetLambdas(s)

	go yatas.CheckTest(checkConfig.Wg, c, "AWS_LMD_001", CheckIfLambdaPrivate)(checkConfig, lambdas, "AWS_LMD_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_LMD_002", CheckIfLambdaInSecurityGroup)(checkConfig, lambdas, "AWS_LMD_002")
	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
