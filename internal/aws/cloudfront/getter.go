package cloudfront

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

type SummaryToConfig struct {
	summary types.DistributionSummary
	config  types.DistributionConfig
}

type CloudfrontGetObjectApi interface {
	GetDistributionConfig(ctx context.Context, params *cloudfront.GetDistributionConfigInput, optFns ...func(*cloudfront.Options)) (*cloudfront.GetDistributionConfigOutput, error)
	ListDistributions(ctx context.Context, params *cloudfront.ListDistributionsInput, optFns ...func(*cloudfront.Options)) (*cloudfront.ListDistributionsOutput, error)
}

func GetAllCloudfront(svc CloudfrontGetObjectApi) []types.DistributionSummary {
	input := &cloudfront.ListDistributionsInput{}
	result, err := svc.ListDistributions(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	return result.DistributionList.Items
}

func GetAllDistributionConfig(svc CloudfrontGetObjectApi, ds []types.DistributionSummary) []SummaryToConfig {
	var d []SummaryToConfig
	for _, cc := range ds {
		input := &cloudfront.GetDistributionConfigInput{
			Id: cc.Id,
		}
		result, err := svc.GetDistributionConfig(context.TODO(), input)
		if err != nil {
			panic(err)
		}
		d = append(d, SummaryToConfig{summary: cc, config: *result.DistributionConfig})
	}
	return d
}
