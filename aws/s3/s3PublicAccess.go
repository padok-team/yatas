package s3

import (
	"fmt"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfS3PublicAccessBlockEnabled(checkConfig yatas.CheckConfig, s3toPublicBlockAccess []S3toPublicBlockAccess, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("S3 bucket have public access block enabled", "Check if S3 buckets are using Public Access Block", testName)
	for _, bucket := range s3toPublicBlockAccess {
		if !bucket.Config {
			Message := "S3 bucket " + bucket.BucketName + " is not using Public Access Block"
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		} else {
			Message := "S3 bucket " + bucket.BucketName + " is using Public Access Block"
			result := yatas.Result{Status: "OK", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
