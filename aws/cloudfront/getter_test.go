package cloudfront

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

type mockGetCloudfront func()

func (m mockGetCloudfront) GetDistributionConfig(ctx context.Context, params *cloudfront.GetDistributionConfigInput, optFns ...func(*cloudfront.Options)) (*cloudfront.GetDistributionConfigOutput, error) {
	return &cloudfront.GetDistributionConfigOutput{
		DistributionConfig: &types.DistributionConfig{
			DefaultCacheBehavior: &types.DefaultCacheBehavior{
				ForwardedValues: &types.ForwardedValues{
					QueryString: aws.Bool(true),
				},
			},
			Enabled: aws.Bool(true),
			Logging: &types.LoggingConfig{
				Enabled: aws.Bool(true),
			},
		},
	}, nil
}

func (m mockGetCloudfront) ListDistributions(ctx context.Context, params *cloudfront.ListDistributionsInput, optFns ...func(*cloudfront.Options)) (*cloudfront.ListDistributionsOutput, error) {
	return &cloudfront.ListDistributionsOutput{
		DistributionList: &types.DistributionList{
			Items: []types.DistributionSummary{
				{
					Id: aws.String("123"),
					DefaultCacheBehavior: &types.DefaultCacheBehavior{
						TargetOriginId: aws.String("123"),
					},
					IsIPV6Enabled: aws.Bool(true),
				},
			},
		},
	}, nil
}

func TestGetAllCloudfront(t *testing.T) {
	type args struct {
		svc CloudfrontGetObjectApi
	}
	tests := []struct {
		name string
		args args
		want []types.DistributionSummary
	}{
		{
			name: "Empty list of Cloudfront distributions",
			args: args{
				svc: mockGetCloudfront(nil),
			},
			want: []types.DistributionSummary{
				{
					Id: aws.String("123"),
					DefaultCacheBehavior: &types.DefaultCacheBehavior{
						TargetOriginId: aws.String("123"),
					},
					IsIPV6Enabled: aws.Bool(true),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAllCloudfront(tt.args.svc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllCloudfront() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestGetAllDistributionConfig(t *testing.T) {
	type args struct {
		svc CloudfrontGetObjectApi
		ds  []types.DistributionSummary
	}
	tests := []struct {
		name string
		args args
		want []SummaryToConfig
	}{
		{
			name: "Empty list of Cloudfront distributions",
			args: args{
				svc: mockGetCloudfront(nil),
				ds:  []types.DistributionSummary{},
			},
			want: []SummaryToConfig{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAllDistributionConfig(tt.args.svc, tt.args.ds); len(got) != len(tt.want) {
				t.Errorf("GetAllDistributionConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
