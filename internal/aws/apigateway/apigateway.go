package apigateway

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func GetApiGateways(s *session.Session) []*apigateway.RestApi {
	svc := apigateway.New(s)
	input := &apigateway.GetRestApisInput{}
	result, err := svc.GetRestApis(input)
	if err != nil {
		return nil
	}
	return result.Items
}

func GetAllResourcesApiGateway(s *session.Session, apiId string) []*apigateway.Resource {
	svc := apigateway.New(s)
	input := &apigateway.GetResourcesInput{
		RestApiId: &apiId,
	}
	result, err := svc.GetResources(input)
	if err != nil {
		return nil
	}
	return result.Items
}

func GetAllStagesApiGateway(s *session.Session, apis []*apigateway.RestApi) []*apigateway.Stage {
	var stages []*apigateway.Stage
	for _, api := range apis {
		svc := apigateway.New(s)
		input := &apigateway.GetStagesInput{
			RestApiId: api.Id,
		}
		result, err := svc.GetStages(input)
		if err != nil {
			return nil
		}
		stages = append(stages, result.Item...)
	}
	return stages
}

func CheckIfStagesCloudwatchLogsExist(s *session.Session, stages []*apigateway.Stage, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "Apigateway Cloudwatch Logs enabled"
	check.Id = testName
	check.Description = "Check if all cloudwatch logs are enabled for all stages"
	check.Status = "OK"
	for _, stage := range stages {
		if stage.AccessLogSettings != nil && stage.AccessLogSettings.DestinationArn != nil {
			check.Status = "OK"
			status := "OK"
			Message := "Cloudwatch logs are enabled on stage" + *stage.StageName
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *stage.StageName})
		} else {
			status := "FAIL"
			Message := "Cloudwatch logs are not enabled on " + *stage.StageName
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *stage.StageName})
		}
	}
	*c = append(*c, check)
}

func CheckIfStagesProtectedByAcl(s *session.Session, stages []*apigateway.Stage, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "Apigateway Stages protected by ACL"
	check.Id = testName
	check.Description = "Check if all stages are protected by ACL"
	check.Status = "OK"
	for _, stage := range stages {
		if *stage.WebAclArn != "" {
			check.Status = "OK"
			status := "OK"
			Message := "Stage " + *stage.StageName + " is protected by ACL"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *stage.StageName})
		} else {
			status := "FAIL"
			Message := "Stage " + *stage.StageName + " is not protected by ACL"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *stage.StageName})
		}
	}
	*c = append(*c, check)
}

func RunApiGatewayTests(s *session.Session, c *config.Config) []types.Check {
	// var checks []types.Check
	var checks []types.Check

	apis := GetApiGateways(s)
	stages := GetAllStagesApiGateway(s, apis)
	config.CheckTest(c, "AWS_APG_001", CheckIfStagesCloudwatchLogsExist)(s, stages, "AWS_APG_001", &checks)
	config.CheckTest(c, "AWS_APG_002", CheckIfStagesProtectedByAcl)(s, stages, "AWS_APG_002", &checks)

	return checks
}
