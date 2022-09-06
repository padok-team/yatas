package lambda

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfLambdaNoErrors(checkConfig yatas.CheckConfig, lambdas []types.FunctionConfiguration, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("Lambdas are not with errors", "Check if all Lambdas are running smoothly", testName)
	for _, lambda := range lambdas {
		if lambda.StateReasonCode != types.StateReasonCodeIdle && lambda.StateReasonCode != "" {
			Message := "Lambda " + *lambda.FunctionName + " is in error with code : " + string(lambda.StateReasonCode)
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: *lambda.FunctionArn}
			check.AddResult(result)
		} else {
			Message := "Lambda " + *lambda.FunctionName + " is running smoothly"
			result := yatas.Result{Status: "OK", Message: Message, ResourceID: *lambda.FunctionArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
