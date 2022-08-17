package apigateway

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
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

func GetAllResourcesApiGateway(checkConfig yatas.CheckConfig, apiId string) []types.Resource {
	svc := apigateway.NewFromConfig(checkConfig.ConfigAWS)
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
