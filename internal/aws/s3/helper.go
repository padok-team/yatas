package s3

import "github.com/aws/aws-sdk-go-v2/service/s3/types"

func OnlyBucketInRegion(BucketAndNotInRegion BucketAndNotInRegion) []types.Bucket {
	var buckets []types.Bucket
	for _, bucket := range BucketAndNotInRegion.Buckets {
		found := false
		for _, bucketNotInRegion := range BucketAndNotInRegion.NotInRegion {
			if *bucket.Name == *bucketNotInRegion.Name {
				found = true
				break
			}
		}
		if !found {
			buckets = append(buckets, bucket)
		}
	}
	return buckets
}
