package lambda

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func GetLambdas(s *session.Session) []*lambda.FunctionConfiguration {
	svc := lambda.New(s)
	input := &lambda.ListFunctionsInput{
		MaxItems: aws.Int64(100),
	}
	result, err := svc.ListFunctions(input)
	if err != nil {
		panic(err)
	}
	return result.Functions
}

func CheckIfLambdaPrivate(s *session.Session, lambdas []*lambda.FunctionConfiguration, c *[]types.Check) {
	logger.Info("Running AWS_LMD_001")
	var check types.Check
	check.Name = "Lambda Private"
	check.Id = "AWS_LMD_001"
	check.Description = "Check if all Lambdas are private"
	check.Status = "OK"
	for _, lambda := range lambdas {
		if lambda.VpcConfig == nil {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Lambda " + *lambda.FunctionName + " is public"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "Lambda " + *lambda.FunctionName + " is private"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func CheckIfLambdaInSecurityGroup(s *session.Session, lambdas []*lambda.FunctionConfiguration, c *[]types.Check) {
	logger.Info("Running AWS_LMD_002")
	var check types.Check
	check.Name = "Lambda In Security Group"
	check.Id = "AWS_LMD_002"
	check.Description = "Check if all Lambdas are in a security group"
	check.Status = "OK"
	for _, lambda := range lambdas {
		if lambda.VpcConfig == nil || lambda.VpcConfig.SecurityGroupIds == nil {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Lambda " + *lambda.FunctionName + " is not in a security group"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "Lambda " + *lambda.FunctionName + " is in a security group"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func RunLambdaTests(s *session.Session) []types.Check {
	var checks []types.Check
	lambdas := GetLambdas(s)
	CheckIfLambdaPrivate(s, lambdas, &checks)
	CheckIfLambdaInSecurityGroup(s, lambdas, &checks)
	return checks
}
