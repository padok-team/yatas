package apigateway

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

type APIGatewayGetObjectAPI interface {
	GetRestApis(ctx context.Context, params *apigateway.GetRestApisInput, optFns ...func(*apigateway.Options)) (*apigateway.GetRestApisOutput, error)
	GetResources(ctx context.Context, params *apigateway.GetResourcesInput, optFns ...func(*apigateway.Options)) (*apigateway.GetResourcesOutput, error)
	GetStages(ctx context.Context, params *apigateway.GetStagesInput, optFns ...func(*apigateway.Options)) (*apigateway.GetStagesOutput, error)
}

func GetApiGateways(svc APIGatewayGetObjectAPI) []types.RestApi {
	input := &apigateway.GetRestApisInput{}
	var apis []types.RestApi
	result, err := svc.GetRestApis(context.TODO(), input)
	apis = append(apis, result.Items...)
	if err != nil {
		return nil
	}
	for {
		if result.Position == nil {
			break
		}
		input.Position = result.Position
		result, err = svc.GetRestApis(context.TODO(), input)
		if err != nil {
			return nil
		}
		apis = append(apis, result.Items...)
	}

	return apis
}

func GetAllResourcesApiGateway(svc APIGatewayGetObjectAPI, apiId string) []types.Resource {
	input := &apigateway.GetResourcesInput{
		RestApiId: &apiId,
	}
	var resources []types.Resource
	result, err := svc.GetResources(context.TODO(), input)
	resources = append(resources, result.Items...)
	if err != nil {
		return nil
	}

	for {
		if result.Position == nil {
			break
		}
		input.Position = result.Position
		result, err = svc.GetResources(context.TODO(), input)
		if err != nil {
			return nil
		}
		resources = append(resources, result.Items...)
	}
	return resources
}

func GetAllStagesApiGateway(svc APIGatewayGetObjectAPI, apis []types.RestApi) map[string][]types.Stage {
	stages := make(map[string][]types.Stage)
	for _, api := range apis {
		input := &apigateway.GetStagesInput{
			RestApiId: api.Id,
		}
		result, err := svc.GetStages(context.TODO(), input)
		if err != nil {
			return nil
		}
		stages[*api.Id] = result.Item

	}
	return stages
}
