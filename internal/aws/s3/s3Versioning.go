package s3

import (
	"fmt"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfBucketObjectVersioningEnabled(checkConfig yatas.CheckConfig, buckets []S3ToVersioning, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("S3 buckets are versioned", "Check if S3 buckets are using object versioning", testName)
	for _, bucket := range buckets {
		if !bucket.Versioning {
			Message := "S3 bucket " + bucket.BucketName + " is not using object versioning"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		} else {
			Message := "S3 bucket " + bucket.BucketName + " is using object versioning"
			result := results.Result{Status: "OK", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		}

	}
	checkConfig.Queue <- check
}
