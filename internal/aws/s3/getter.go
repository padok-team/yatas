package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stangirard/yatas/internal/logger"
)

func GetListS3(s aws.Config) []types.Bucket {
	logger.Debug("Getting list of S3 buckets")
	svc := s3.NewFromConfig(s)

	params := &s3.ListBucketsInput{}
	resp, err := svc.ListBuckets(context.TODO(), params)
	if err != nil {
		panic(err)
	}

	logger.Debug(fmt.Sprintf("%v", resp.Buckets))
	return resp.Buckets
}

func GetListS3NotInRegion(s aws.Config, region string) []types.Bucket {
	logger.Debug("Getting list of S3 buckets not in region")
	svc := s3.NewFromConfig(s)

	params := &s3.ListBucketsInput{}
	resp, err := svc.ListBuckets(context.TODO(), params)
	if err != nil {
		panic(err)
	}

	var buckets []types.Bucket
	for _, bucket := range resp.Buckets {
		if !CheckS3Location(s, *bucket.Name, region) {
			buckets = append(buckets, bucket)
		}
	}
	logger.Debug(fmt.Sprintf("%v", buckets))
	return buckets
}
