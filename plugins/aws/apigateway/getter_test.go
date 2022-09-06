package apigateway

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

type mockGetApiGateways func()

func (m mockGetApiGateways) GetRestApis(ctx context.Context, input *apigateway.GetRestApisInput, optFns ...func(*apigateway.Options)) (*apigateway.GetRestApisOutput, error) {
	// Return an empty list of API Gateway instances
	timeTest, _ := time.Parse("2006-01-02T15:04:05Z", "2019-01-01T00:00:00Z")
	return &apigateway.GetRestApisOutput{
		Items: []types.RestApi{
			{
				Id:                     aws.String("id"),
				CreatedDate:            aws.Time(timeTest),
				Name:                   aws.String("name"),
				MinimumCompressionSize: aws.Int32(0),
				Version:                aws.String("version"),
			},
		},
	}, nil
}

func (m mockGetApiGateways) GetResources(ctx context.Context, input *apigateway.GetResourcesInput, optFns ...func(*apigateway.Options)) (*apigateway.GetResourcesOutput, error) {
	// Return an empty list of API Gateway resources
	return &apigateway.GetResourcesOutput{
		Items: []types.Resource{
			{
				Path:     aws.String("path"),
				Id:       aws.String("id"),
				ParentId: aws.String("parentId"),
			},
		},
	}, nil
}

func (m mockGetApiGateways) GetStages(ctx context.Context, input *apigateway.GetStagesInput, optFns ...func(*apigateway.Options)) (*apigateway.GetStagesOutput, error) {
	// Return an empty list of API Gateway stages
	return &apigateway.GetStagesOutput{
		Item: []types.Stage{
			{
				DeploymentId: aws.String("deploymentId"),
				AccessLogSettings: &types.AccessLogSettings{
					DestinationArn: aws.String("destinationArn"),
					Format:         aws.String("format"),
				},
				TracingEnabled: true,
				WebAclArn:      aws.String("webAclArn"),
			},
		},
	}, nil
}

func TestGetApiGateways(t *testing.T) {
	type args struct {
		svc APIGatewayGetObjectAPI
	}

	timeTest, _ := time.Parse("2006-01-02T15:04:05Z", "2019-01-01T00:00:00Z")
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
			want: []types.RestApi{
				{
					Id:                     aws.String("id"),
					CreatedDate:            aws.Time(timeTest),
					Name:                   aws.String("name"),
					MinimumCompressionSize: aws.Int32(0),
					Version:                aws.String("version"),
				},
			},
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
			want: []types.Resource{
				{
					Path:     aws.String("path"),
					Id:       aws.String("id"),
					ParentId: aws.String("parentId"),
				},
			},
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
		want map[string][]types.Stage
	}{
		{
			name: "Empty list of API Gateway stages",
			args: args{
				svc: mockGetApiGateways(nil),
				apis: []types.RestApi{
					{
						Id: aws.String("test"),
					},
				},
			},
			want: map[string][]types.Stage{
				"test": {

					{
						DeploymentId: aws.String("deploymentId"),
						AccessLogSettings: &types.AccessLogSettings{
							DestinationArn: aws.String("destinationArn"),
							Format:         aws.String("format"),
						},
						TracingEnabled: true,
						WebAclArn:      aws.String("webAclArn"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAllStagesApiGateway(tt.args.svc, tt.args.apis); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllStagesApiGateway() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
