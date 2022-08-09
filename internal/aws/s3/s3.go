package s3

import (
	"context"
	"fmt"
	"strings"

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

func checkIfEncryptionEnabled(s aws.Config, buckets []types.Bucket, testName string, c *[]results.Check) {
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
}

func CheckIfBucketInOneZone(s aws.Config, buckets []types.Bucket, testName string, c *[]results.Check) {
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
}

func CheckIfBucketObjectVersioningEnabled(s aws.Config, buckets []types.Bucket, testName string, c *[]results.Check) {
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
}

func CheckIfObjectLockConfigurationEnabled(s aws.Config, buckets []types.Bucket, testName string, c *[]results.Check) {
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

func RunS3Test(s aws.Config, c *yatas.Config) []results.Check {
	var checks []results.Check
	logger.Debug("Starting S3 tests")
	buckets := GetListS3(s)
	yatas.CheckTest(c, "AWS_S3_001", checkIfEncryptionEnabled)(s, buckets, "AWS_S3_001", &checks)
	yatas.CheckTest(c, "AWS_S3_002", CheckIfBucketInOneZone)(s, buckets, "AWS_S3_002", &checks)
	yatas.CheckTest(c, "AWS_S3_003", CheckIfBucketObjectVersioningEnabled)(s, buckets, "AWS_S3_003", &checks)
	yatas.CheckTest(c, "AWS_S3_004", CheckIfObjectLockConfigurationEnabled)(s, buckets, "AWS_S3_004", &checks)
	return checks
}
