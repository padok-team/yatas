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

func checkIfEncryptionEnabled(wg *sync.WaitGroup, s aws.Config, buckets []types.Bucket, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("S3 Encryption", "Check if S3 encryption is enabled", testName)
	svc := s3.NewFromConfig(s)
	for _, bucket := range buckets {
		if !CheckS3Location(s, *bucket.Name, s.Region) {
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
	queueToAdd <- check
}

func CheckIfBucketInOneZone(wg *sync.WaitGroup, s aws.Config, buckets []types.Bucket, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("S3 Bucket in one zone", "Check if S3 buckets are in one zone", testName)
	for _, bucket := range buckets {
		if !CheckS3Location(s, *bucket.Name, s.Region) {
			Message := "S3 bucket " + *bucket.Name + " is not in the same zone as the account"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *bucket.Name}
			check.AddResult(result)
		} else {
			Message := "S3 bucket " + *bucket.Name + " is in the same zone as the account"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *bucket.Name}
			check.AddResult(result)
		}
	}
	queueToAdd <- check
}

func CheckIfBucketObjectVersioningEnabled(wg *sync.WaitGroup, s aws.Config, buckets []types.Bucket, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("S3 Bucket object versioning", "Check if S3 buckets are using object versioning", testName)
	svc := s3.NewFromConfig(s)
	for _, bucket := range buckets {
		if !CheckS3Location(s, *bucket.Name, s.Region) {
			continue
		}
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
	queueToAdd <- check
}

func CheckIfObjectLockConfigurationEnabled(wg *sync.WaitGroup, s aws.Config, buckets []types.Bucket, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("S3 Bucket retention policy", "Check if S3 buckets are using retention policy", testName)
	svc := s3.NewFromConfig(s)
	for _, bucket := range buckets {
		if !CheckS3Location(s, *bucket.Name, s.Region) {
			continue
		}
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
	queueToAdd <- check
}

func CheckIfS3PublicAccessBlockEnabled(wg *sync.WaitGroup, s aws.Config, buckets []types.Bucket, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("S3 Public Access Block", "Check if S3 buckets are using Public Access Block", testName)
	svc := s3.NewFromConfig(s)
	for _, bucket := range buckets {
		if !CheckS3Location(s, *bucket.Name, s.Region) {
			continue
		}
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
	queueToAdd <- check
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

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []results.Check
	buckets := GetListS3(s)

	go yatas.CheckTest(checkConfig.Wg, c, "AWS_S3_001", checkIfEncryptionEnabled)(checkConfig.Wg, checkConfig.ConfigAWS, buckets, "AWS_S3_001", checkConfig.Queue)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_S3_002", CheckIfBucketInOneZone)(checkConfig.Wg, checkConfig.ConfigAWS, buckets, "AWS_S3_002", checkConfig.Queue)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_S3_003", CheckIfBucketObjectVersioningEnabled)(checkConfig.Wg, checkConfig.ConfigAWS, buckets, "AWS_S3_003", checkConfig.Queue)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_S3_004", CheckIfObjectLockConfigurationEnabled)(checkConfig.Wg, checkConfig.ConfigAWS, buckets, "AWS_S3_004", checkConfig.Queue)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_S3_005", CheckIfS3PublicAccessBlockEnabled)(checkConfig.Wg, checkConfig.ConfigAWS, buckets, "AWS_S3_005", checkConfig.Queue)
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
