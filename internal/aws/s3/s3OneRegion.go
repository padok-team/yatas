package s3

import (
	"fmt"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfBucketInOneZone(checkConfig yatas.CheckConfig, buckets BucketAndNotInRegion, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("S3 buckets are not global but in one zone", "Check if S3 buckets are in one zone", testName)
	for _, bucket := range buckets.Buckets {
		found := false
		for _, region := range buckets.NotInRegion {
			if *bucket.Name == *region.Name {
				Message := "S3 bucket " + *bucket.Name + " but should be in " + checkConfig.ConfigAWS.Region
				result := results.Result{Status: "FAIL", Message: Message, ResourceID: *bucket.Name}
				check.AddResult(result)
				found = true
				break
			}
		}
		if !found {
			Message := "S3 bucket " + *bucket.Name + " is in " + checkConfig.ConfigAWS.Region
			result := results.Result{Status: "OK", Message: Message, ResourceID: *bucket.Name}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
