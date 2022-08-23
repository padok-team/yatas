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
	result, err := svc.GetRestApis(context.TODO(), input)
	if err != nil {
		return nil
	}

	return result.Items
}

func GetAllResourcesApiGateway(svc APIGatewayGetObjectAPI, apiId string) []types.Resource {
	input := &apigateway.GetResourcesInput{
		RestApiId: &apiId,
	}
	result, err := svc.GetResources(context.TODO(), input)
	if err != nil {
		return nil
	}
	return result.Items
}

func GetAllStagesApiGateway(svc APIGatewayGetObjectAPI, apis []types.RestApi) []types.Stage {
	var stages []types.Stage
	for _, api := range apis {
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
