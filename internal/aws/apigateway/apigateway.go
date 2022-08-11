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
	check.Name = "Apigateway Cloudwatch Logs enabled"
	check.Id = testName
	check.Description = "Check if all cloudwatch logs are enabled for all stages"
	check.Status = "OK"
	for _, stage := range stages {
		if stage.AccessLogSettings != nil && stage.AccessLogSettings.DestinationArn != nil {
			check.Status = "OK"
			status := "OK"
			Message := "Cloudwatch logs are enabled on stage" + *stage.StageName
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *stage.StageName})
		} else {
			status := "FAIL"
			Message := "Cloudwatch logs are not enabled on " + *stage.StageName
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *stage.StageName})
		}
	}
	queueToAdd <- check
}

func CheckIfStagesProtectedByAcl(wg *sync.WaitGroup, s aws.Config, stages []types.Stage, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "Apigateway Stages protected by ACL"
	check.Id = testName
	check.Description = "Check if all stages are protected by ACL"
	check.Status = "OK"
	for _, stage := range stages {
		if *stage.WebAclArn != "" {
			check.Status = "OK"
			status := "OK"
			Message := "Stage " + *stage.StageName + " is protected by ACL"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *stage.StageName})
		} else {
			status := "FAIL"
			Message := "Stage " + *stage.StageName + " is not protected by ACL"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *stage.StageName})
		}
	}
	queueToAdd <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	// var checks []results.Check
	var checks []results.Check
	var wg sync.WaitGroup
	queueResults := make(chan results.Check, 10)
	apis := GetApiGateways(s)
	stages := GetAllStagesApiGateway(s, apis)
	go yatas.CheckTest(&wg, c, "AWS_APG_001", CheckIfStagesCloudwatchLogsExist)(&wg, s, stages, "AWS_APG_001", queueResults)
	go yatas.CheckTest(&wg, c, "AWS_APG_002", CheckIfStagesProtectedByAcl)(&wg, s, stages, "AWS_APG_002", queueResults)

	go func() {
		for t := range queueResults {
			checks = append(checks, t)
			wg.Done()
		}
	}()

	wg.Wait()

	queue <- checks
}
