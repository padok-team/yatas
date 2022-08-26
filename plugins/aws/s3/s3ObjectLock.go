package s3

import (
	"fmt"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfObjectLockConfigurationEnabled(checkConfig yatas.CheckConfig, buckets []S3ToObjectLock, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("S3 buckets have a retention policy", "Check if S3 buckets are using retention policy", testName)
	for _, bucket := range buckets {
		if !bucket.ObjectLock {
			Message := "S3 bucket " + bucket.BucketName + " is not using retention policy"
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		} else {
			Message := "S3 bucket " + bucket.BucketName + " is using retention policy"
			result := yatas.Result{Status: "OK", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		}

	}
	checkConfig.Queue <- check
}
