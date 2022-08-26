package s3

import (
	"fmt"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func checkIfEncryptionEnabled(checkConfig yatas.CheckConfig, buckets []S3ToEncryption, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("S3 are encrypted", "Check if S3 encryption is enabled", testName)
	for _, bucket := range buckets {
		if !bucket.Encrypted {
			Message := "S3 bucket " + bucket.BucketName + " is not using encryption"
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		} else {
			Message := "S3 bucket " + bucket.BucketName + " is using encryption"
			result := yatas.Result{Status: "OK", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		}

	}
	checkConfig.Queue <- check
}
