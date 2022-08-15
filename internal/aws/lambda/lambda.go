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

func CheckIfLambdaPrivate(wg *sync.WaitGroup, s aws.Config, lambdas []types.FunctionConfiguration, testName string, queueToAdd chan results.Check) {
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
	queueToAdd <- check
}

func CheckIfLambdaInSecurityGroup(wg *sync.WaitGroup, s aws.Config, lambdas []types.FunctionConfiguration, testName string, queueToAdd chan results.Check) {
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
	queueToAdd <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checks []results.Check
	lambdas := GetLambdas(s)
	var wg sync.WaitGroup
	queueResults := make(chan results.Check, 10)

	go yatas.CheckTest(&wg, c, "AWS_LMD_001", CheckIfLambdaPrivate)(&wg, s, lambdas, "AWS_LMD_001", queueResults)
	go yatas.CheckTest(&wg, c, "AWS_LMD_002", CheckIfLambdaInSecurityGroup)(&wg, s, lambdas, "AWS_LMD_002", queueResults)
	go func() {
		for t := range queueResults {
			checks = append(checks, t)
			wg.Done()
		}
	}()

	wg.Wait()

	queue <- checks
}
