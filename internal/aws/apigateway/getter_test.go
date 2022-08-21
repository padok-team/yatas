package apigateway

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

type mockGetApiGateways func()

func (m mockGetApiGateways) GetRestApis(ctx context.Context, input *apigateway.GetRestApisInput, optFns ...func(*apigateway.Options)) (*apigateway.GetRestApisOutput, error) {
	// Return an empty list of API Gateway instances
	return &apigateway.GetRestApisOutput{
		Items: []types.RestApi{},
	}, nil
}

func (m mockGetApiGateways) GetResources(ctx context.Context, input *apigateway.GetResourcesInput, optFns ...func(*apigateway.Options)) (*apigateway.GetResourcesOutput, error) {
	// Return an empty list of API Gateway resources
	return &apigateway.GetResourcesOutput{
		Items: []types.Resource{},
	}, nil
}

func (m mockGetApiGateways) GetStages(ctx context.Context, input *apigateway.GetStagesInput, optFns ...func(*apigateway.Options)) (*apigateway.GetStagesOutput, error) {
	// Return an empty list of API Gateway stages
	return &apigateway.GetStagesOutput{
		Item: []types.Stage{},
	}, nil
}

func TestGetApiGateways(t *testing.T) {
	type args struct {
		svc APIGatewayGetObjectAPI
	}
	tests := []struct {
		name string
		args args
		want []types.RestApi
	}{
		{
			name: "Empty list of API Gateway instances",
			args: args{
				svc: mockGetApiGateways(nil),
			},
			want: []types.RestApi{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetApiGateways(tt.args.svc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetApiGateways() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllResourcesApiGateway(t *testing.T) {
	type args struct {
		svc   APIGatewayGetObjectAPI
		apiId string
	}
	tests := []struct {
		name string
		args args
		want []types.Resource
	}{
		{
			name: "Empty list of API Gateway resources",
			args: args{
				svc:   mockGetApiGateways(nil),
				apiId: "",
			},
			want: []types.Resource{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAllResourcesApiGateway(tt.args.svc, tt.args.apiId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllResourcesApiGateway() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllStagesApiGateway(t *testing.T) {
	type args struct {
		svc  APIGatewayGetObjectAPI
		apis []types.RestApi
	}
	tests := []struct {
		name string
		args args
		want []types.Stage
	}{
		{
			name: "Empty list of API Gateway stages",
			args: args{
				svc: mockGetApiGateways(nil),
				apis: []types.RestApi{
					{
						Id: aws.String(""),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAllStagesApiGateway(tt.args.svc, tt.args.apis); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllStagesApiGateway() = %v, want %v", got, tt.want)
			}
		})
	}
}
