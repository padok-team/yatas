package s3

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func checkIfEncryptionEnabled(checkConfig yatas.CheckConfig, buckets []S3ToEncryption, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("S3 Encryption", "Check if S3 encryption is enabled", testName)
	for _, bucket := range buckets {
		if !bucket.Encrypted {
			Message := "S3 bucket " + bucket.BucketName + " is not using encryption"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		} else {
			Message := "S3 bucket " + bucket.BucketName + " is using encryption"
			result := results.Result{Status: "OK", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		}

	}
	checkConfig.Queue <- check
}

func CheckIfObjectLockConfigurationEnabled(checkConfig yatas.CheckConfig, buckets []types.Bucket, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("S3 Bucket retention policy", "Check if S3 buckets are using retention policy", testName)
	svc := s3.NewFromConfig(checkConfig.ConfigAWS)
	for _, bucket := range buckets {
		params := &s3.GetObjectLockConfigurationInput{
			Bucket: aws.String(*bucket.Name),
		}
		resp, err := svc.GetObjectLockConfiguration(context.TODO(), params)
		if err != nil || (resp.ObjectLockConfiguration != nil && resp.ObjectLockConfiguration.ObjectLockEnabled != "Enabled") {
			Message := "S3 bucket " + *bucket.Name + " is not using retention policy"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *bucket.Name}
			check.AddResult(result)

		} else {
			Message := "S3 bucket " + *bucket.Name + " is using retention policy"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *bucket.Name}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

func CheckS3Location(s aws.Config, bucket, region string) bool {
	logger.Debug("Getting S3 location")
	svc := s3.NewFromConfig(s)

	params := &s3.GetBucketLocationInput{
		Bucket: aws.String(bucket),
	}
	resp, err := svc.GetBucketLocation(context.TODO(), params)
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
		return false
	}
	logger.Debug(fmt.Sprintf("%v", resp))

	if resp.LocationConstraint != "" {
		if string(resp.LocationConstraint) == region {
			return true
		} else {
			return false
		}

	} else {
		return false
	}
}

type BucketAndNotInRegion struct {
	Buckets     []types.Bucket
	NotInRegion []types.Bucket
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []results.Check
	buckets := GetListS3(s)
	bucketsNotInRegion := GetListS3NotInRegion(s, s.Region)
	couple := BucketAndNotInRegion{buckets, bucketsNotInRegion}
	OnlyBucketInRegion := OnlyBucketInRegion(couple)
	S3ToEncryption := GetS3ToEncryption(s, OnlyBucketInRegion)
	S3toPublicBlockAccess := GetS3ToPublicBlockAccess(s, OnlyBucketInRegion)
	S3ToVersioning := GetS3ToVersioning(s, OnlyBucketInRegion)

	go yatas.CheckTest(checkConfig.Wg, c, "AWS_S3_001", checkIfEncryptionEnabled)(checkConfig, S3ToEncryption, "AWS_S3_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_S3_002", CheckIfBucketInOneZone)(checkConfig, couple, "AWS_S3_002")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_S3_003", CheckIfBucketObjectVersioningEnabled)(checkConfig, S3ToVersioning, "AWS_S3_003")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_S3_004", CheckIfObjectLockConfigurationEnabled)(checkConfig, OnlyBucketInRegion, "AWS_S3_004")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_S3_005", CheckIfS3PublicAccessBlockEnabled)(checkConfig, S3toPublicBlockAccess, "AWS_S3_005")
	// Wait for all the goroutines to finish

	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
