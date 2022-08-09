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

func checkIfEncryptionEnabled(wg *sync.WaitGroup, s aws.Config, buckets []types.Bucket, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "S3 Encryption"
	check.Id = testName
	check.Description = "Check if S3 encryption is enabled"
	check.Status = "OK"
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
			check.Status = "FAIL"
			status := "FAIL"
			Message := "S3 encryption is not enabled on " + *bucket.Name
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *bucket.Name})
		} else {
			status := "OK"
			Message := "S3 encryption is enabled on " + *bucket.Name
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *bucket.Name})
		}
	}
	*c = append(*c, check)
	wg.Done()
}

func CheckIfBucketInOneZone(wg *sync.WaitGroup, s aws.Config, buckets []types.Bucket, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "S3 Bucket in one zone"
	check.Id = testName
	check.Description = "Check if S3 buckets are in one zone"
	check.Status = "OK"
	for _, bucket := range buckets {
		if !CheckS3Location(s, *bucket.Name, s.Region) {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "S3 bucket " + *bucket.Name + " is not in the same zone as the account"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *bucket.Name})
		} else {
			status := "OK"
			Message := "S3 bucket " + *bucket.Name + " is in the same zone as the account"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *bucket.Name})
		}
	}
	*c = append(*c, check)
	wg.Done()
}

func CheckIfBucketObjectVersioningEnabled(wg *sync.WaitGroup, s aws.Config, buckets []types.Bucket, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "S3 Bucket object versioning"
	check.Id = testName
	check.Description = "Check if S3 buckets are using object versioning"
	check.Status = "OK"
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
		if &resp.Status != nil && resp.Status != "Enabled" {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "S3 bucket " + *bucket.Name + " is not using object versioning"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *bucket.Name})
		} else {
			status := "OK"
			Message := "S3 bucket " + *bucket.Name + " is using object versioning"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *bucket.Name})
		}
	}
	*c = append(*c, check)
	wg.Done()
}

func CheckIfObjectLockConfigurationEnabled(wg *sync.WaitGroup, s aws.Config, buckets []types.Bucket, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "S3 Bucket retention policy"
	check.Id = testName
	check.Description = "Check if S3 buckets are using retention policy"
	check.Status = "OK"
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
			check.Status = "FAIL"
			status := "FAIL"
			Message := "S3 bucket " + *bucket.Name + " is not using retention policy"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *bucket.Name})

		} else {
			status := "OK"
			Message := "S3 bucket " + *bucket.Name + " is using retention policy"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *bucket.Name})
		}
	}
	*c = append(*c, check)
	wg.Done()
}

func CheckIfS3PublicAccessBlockEnabled(wg *sync.WaitGroup, s aws.Config, buckets []types.Bucket, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "S3 Public Access Block"
	check.Id = testName
	check.Description = "Check if S3 buckets are using Public Access Block"
	check.Status = "OK"
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
			check.Status = "FAIL"
			status := "FAIL"
			Message := "S3 bucket " + *bucket.Name + " is not using Public Access Block"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *bucket.Name})
		} else {
			status := "OK"
			Message := "S3 bucket " + *bucket.Name + " is using Public Access Block"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *bucket.Name})
		}
	}
	*c = append(*c, check)
	wg.Done()
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

	var checks []results.Check
	logger.Debug("Starting S3 tests")
	buckets := GetListS3(s)
	// Create a wait group to wait for all the goroutines to finish
	var wg sync.WaitGroup
	go yatas.CheckTest(&wg, c, "AWS_S3_001", checkIfEncryptionEnabled)(&wg, s, buckets, "AWS_S3_001", &checks)
	go yatas.CheckTest(&wg, c, "AWS_S3_002", CheckIfBucketInOneZone)(&wg, s, buckets, "AWS_S3_002", &checks)
	go yatas.CheckTest(&wg, c, "AWS_S3_003", CheckIfBucketObjectVersioningEnabled)(&wg, s, buckets, "AWS_S3_003", &checks)
	go yatas.CheckTest(&wg, c, "AWS_S3_004", CheckIfObjectLockConfigurationEnabled)(&wg, s, buckets, "AWS_S3_004", &checks)
	go yatas.CheckTest(&wg, c, "AWS_S3_005", CheckIfS3PublicAccessBlockEnabled)(&wg, s, buckets, "AWS_S3_005", &checks)
	// Wait for all the goroutines to finish
	wg.Wait()

	queue <- checks
}
