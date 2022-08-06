package s3

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func GetListS3(s *session.Session) []*s3.Bucket {
	logger.Debug("Getting list of S3 buckets")
	svc := s3.New(s)

	params := &s3.ListBucketsInput{}
	resp, err := svc.ListBuckets(params)
	if err != nil {
		panic(err)
	}

	logger.Debug(fmt.Sprintf("%v", resp.Buckets))
	return resp.Buckets
}

func checkIfEncryptionEnabled(s *session.Session, buckets []*s3.Bucket, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "S3 Encryption"
	check.Id = testName
	check.Description = "Check if S3 encryption is enabled"
	check.Status = "OK"
	svc := s3.New(s)
	for _, bucket := range buckets {
		if !CheckS3Location(s, *bucket.Name, *s.Config.Region) {
			continue
		}
		params := &s3.GetBucketEncryptionInput{
			Bucket: aws.String(*bucket.Name),
		}
		_, err := svc.GetBucketEncryption(params)
		// If error contains ServerSideEncryptionConfigurationNotFoundError, then err is nil
		if err != nil && !strings.Contains(err.Error(), "ServerSideEncryptionConfigurationNotFoundError") {
			panic(err)
		} else if err != nil {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "S3 encryption is not enabled on " + *bucket.Name
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "S3 encryption is enabled on " + *bucket.Name
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func CheckIfBucketInOneZone(s *session.Session, buckets []*s3.Bucket, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "S3 Bucket in one zone"
	check.Id = testName
	check.Description = "Check if S3 buckets are in one zone"
	check.Status = "OK"
	for _, bucket := range buckets {
		if !CheckS3Location(s, *bucket.Name, *s.Config.Region) {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "S3 bucket " + *bucket.Name + " is not in the same zone as the account"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "S3 bucket " + *bucket.Name + " is in the same zone as the account"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func CheckS3Location(s *session.Session, bucket, region string) bool {
	logger.Debug("Getting S3 location")
	svc := s3.New(s)

	params := &s3.GetBucketLocationInput{
		Bucket: aws.String(bucket),
	}
	resp, err := svc.GetBucketLocation(params)
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
		return false
	}
	logger.Debug(fmt.Sprintf("%v", resp))

	if resp.LocationConstraint == nil {
		return false
	} else if *resp.LocationConstraint != "" {
		if *resp.LocationConstraint == region {
			return true
		} else {
			return false
		}

	} else {
		return false
	}
}

type testFunc func(*session.Session, []*s3.Bucket, []types.Check)

func RunS3Test(s *session.Session, c *config.Config) []types.Check {
	var checks []types.Check
	logger.Debug("Starting S3 tests")
	buckets := GetListS3(s)
	config.CheckTest(c, "AWS_S3_001", checkIfEncryptionEnabled)(s, buckets, "AWS_S3_001", &checks)
	config.CheckTest(c, "AWS_S3_002", CheckIfBucketInOneZone)(s, buckets, "AWS_S3_002", &checks)
	return checks
}
