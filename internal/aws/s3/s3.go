package s3

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func checkIfEncryptionEnabled(checkConfig yatas.CheckConfig, buckets []types.Bucket, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("S3 Encryption", "Check if S3 encryption is enabled", testName)
	svc := s3.NewFromConfig(checkConfig.ConfigAWS)
	for _, bucket := range buckets {
		if !CheckS3Location(checkConfig.ConfigAWS, *bucket.Name, checkConfig.ConfigAWS.Region) {
			continue
		}
		params := &s3.GetBucketEncryptionInput{
			Bucket: aws.String(*bucket.Name),
		}
		_, err := svc.GetBucketEncryption(context.TODO(), params)
		// If error contains ServerSideEncryptionConfigurationNotFoundError, then err is nil
		if err != nil && !strings.Contains(err.Error(), "ServerSideEncryptionConfigurationNotFoundError") {
			panic(err)
		} else if err != nil {
			Message := "S3 encryption is not enabled on " + *bucket.Name
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *bucket.Name}
			check.AddResult(result)
		} else {
			Message := "S3 encryption is enabled on " + *bucket.Name
			result := results.Result{Status: "OK", Message: Message, ResourceID: *bucket.Name}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

func CheckIfBucketObjectVersioningEnabled(checkConfig yatas.CheckConfig, buckets []types.Bucket, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("S3 Bucket object versioning", "Check if S3 buckets are using object versioning", testName)
	svc := s3.NewFromConfig(checkConfig.ConfigAWS)
	for _, bucket := range buckets {
		params := &s3.GetBucketVersioningInput{
			Bucket: aws.String(*bucket.Name),
		}
		resp, err := svc.GetBucketVersioning(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if resp.Status != "Enabled" {
			Message := "S3 bucket " + *bucket.Name + " is not using object versioning"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *bucket.Name}
			check.AddResult(result)
		} else {
			Message := "S3 bucket " + *bucket.Name + " is using object versioning"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *bucket.Name}
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

func CheckIfS3PublicAccessBlockEnabled(checkConfig yatas.CheckConfig, buckets []types.Bucket, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("S3 Public Access Block", "Check if S3 buckets are using Public Access Block", testName)
	svc := s3.NewFromConfig(checkConfig.ConfigAWS)
	for _, bucket := range buckets {
		params := &s3.GetPublicAccessBlockInput{
			Bucket: aws.String(*bucket.Name),
		}
		resp, err := svc.GetPublicAccessBlock(context.TODO(), params)

		if err != nil || resp.PublicAccessBlockConfiguration == nil || !resp.PublicAccessBlockConfiguration.BlockPublicAcls {
			Message := "S3 bucket " + *bucket.Name + " is not using Public Access Block"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *bucket.Name}
			check.AddResult(result)
		} else {
			Message := "S3 bucket " + *bucket.Name + " is using Public Access Block"
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

	go yatas.CheckTest(checkConfig.Wg, c, "AWS_S3_001", checkIfEncryptionEnabled)(checkConfig, OnlyBucketInRegion, "AWS_S3_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_S3_002", CheckIfBucketInOneZone)(checkConfig, couple, "AWS_S3_002")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_S3_003", CheckIfBucketObjectVersioningEnabled)(checkConfig, OnlyBucketInRegion, "AWS_S3_003")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_S3_004", CheckIfObjectLockConfigurationEnabled)(checkConfig, OnlyBucketInRegion, "AWS_S3_004")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_S3_005", CheckIfS3PublicAccessBlockEnabled)(checkConfig, OnlyBucketInRegion, "AWS_S3_005")
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
