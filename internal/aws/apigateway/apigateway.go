package apigateway

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func GetApiGateways(s aws.Config) []types.RestApi {
	svc := apigateway.NewFromConfig(s)
	input := &apigateway.GetRestApisInput{}
	result, err := svc.GetRestApis(context.TODO(), input)
	if err != nil {
		return nil
	}
	return result.Items
}

func GetAllResourcesApiGateway(wg *sync.WaitGroup, s aws.Config, apiId string) []types.Resource {
	svc := apigateway.NewFromConfig(s)
	input := &apigateway.GetResourcesInput{
		RestApiId: &apiId,
	}
	result, err := svc.GetResources(context.TODO(), input)
	if err != nil {
		return nil
	}
	return result.Items
}

func GetAllStagesApiGateway(s aws.Config, apis []types.RestApi) []types.Stage {
	var stages []types.Stage
	for _, api := range apis {
		svc := apigateway.NewFromConfig(s)
		input := &apigateway.GetStagesInput{
			RestApiId: api.Id,
		}
		result, err := svc.GetStages(context.TODO(), input)
		if err != nil {
			return nil
		}
		stages = append(stages, result.Item...)
	}
	return stages
}

func CheckIfStagesCloudwatchLogsExist(wg *sync.WaitGroup, s aws.Config, stages []types.Stage, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Apigateway Cloudwatch Logs enabled", "Check if all cloudwatch logs are enabled for all stages", testName)
	for _, stage := range stages {
		if stage.AccessLogSettings != nil && stage.AccessLogSettings.DestinationArn != nil {
			Message := "Cloudwatch logs are enabled on stage" + *stage.StageName
			result := results.Result{Status: "OK", Message: Message, ResourceID: *stage.StageName}
			check.AddResult(result)
		} else {
			Message := "Cloudwatch logs are not enabled on " + *stage.StageName
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *stage.StageName}
			check.AddResult(result)
		}
	}
	queueToAdd <- check
}

func CheckIfStagesProtectedByAcl(wg *sync.WaitGroup, s aws.Config, stages []types.Stage, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("APIGateway stages protected, by ACL", "Check if all stages are protected by ACL", testName)
	for _, stage := range stages {
		if *stage.WebAclArn != "" {
			Message := "Stage " + *stage.StageName + " is protected by ACL"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *stage.StageName}
			check.AddResult(result)
		} else {
			Message := "Stage " + *stage.StageName + " is not protected by ACL"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *stage.StageName}
			check.AddResult(result)
		}
	}
	queueToAdd <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {
	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []results.Check
	apis := GetApiGateways(s)
	stages := GetAllStagesApiGateway(s, apis)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_APG_001", CheckIfStagesCloudwatchLogsExist)(checkConfig.Wg, checkConfig.ConfigAWS, stages, "AWS_APG_001", checkConfig.Queue)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_APG_002", CheckIfStagesProtectedByAcl)(checkConfig.Wg, checkConfig.ConfigAWS, stages, "AWS_APG_002", checkConfig.Queue)

	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
