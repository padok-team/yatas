package cloudfront

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func GetAllCloudfront(s aws.Config) []types.DistributionSummary {
	svc := cloudfront.NewFromConfig(s)
	input := &cloudfront.ListDistributionsInput{}
	result, err := svc.ListDistributions(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	return result.DistributionList.Items
}

func GetAllDistributionConfig(s aws.Config, ds []types.DistributionSummary) []types.DistributionConfig {
	svc := cloudfront.NewFromConfig(s)
	var d []types.DistributionConfig
	for _, cc := range ds {
		input := &cloudfront.GetDistributionConfigInput{
			Id: cc.Id,
		}
		result, err := svc.GetDistributionConfig(context.TODO(), input)
		if err != nil {
			panic(err)
		}
		d = append(d, *result.DistributionConfig)
	}
	return d
}
