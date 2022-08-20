package s3

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stangirard/yatas/internal/logger"
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

func GetListS3NotInRegion(s aws.Config, region string) []types.Bucket {
	logger.Debug("Getting list of S3 buckets not in region")
	svc := s3.NewFromConfig(s)

	params := &s3.ListBucketsInput{}
	resp, err := svc.ListBuckets(context.TODO(), params)
	if err != nil {
		panic(err)
	}

	var buckets []types.Bucket
	for _, bucket := range resp.Buckets {
		if !CheckS3Location(s, *bucket.Name, region) {
			buckets = append(buckets, bucket)
		}
	}
	logger.Debug(fmt.Sprintf("%v", buckets))
	return buckets
}

type S3toPublicBlockAccess struct {
	BucketName string
	Config     bool
}

func GetS3ToPublicBlockAccess(s aws.Config, b []types.Bucket) []S3toPublicBlockAccess {
	logger.Debug("Getting list of S3 buckets not in region")
	svc := s3.NewFromConfig(s)

	var s3toPublicBlockAccess []S3toPublicBlockAccess
	for _, bucket := range b {
		params := &s3.GetPublicAccessBlockInput{
			Bucket: aws.String(*bucket.Name),
		}
		resp, err := svc.GetPublicAccessBlock(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if resp.PublicAccessBlockConfiguration != nil && resp.PublicAccessBlockConfiguration.BlockPublicAcls {
			s3toPublicBlockAccess = append(s3toPublicBlockAccess, S3toPublicBlockAccess{*bucket.Name, true})
		} else {
			s3toPublicBlockAccess = append(s3toPublicBlockAccess, S3toPublicBlockAccess{*bucket.Name, false})
		}
	}
	logger.Debug(fmt.Sprintf("%v", s3toPublicBlockAccess))
	return s3toPublicBlockAccess
}

type S3ToEncryption struct {
	BucketName string
	Encrypted  bool
}

func GetS3ToEncryption(s aws.Config, b []types.Bucket) []S3ToEncryption {
	logger.Debug("Getting list of S3 buckets not in region")
	svc := s3.NewFromConfig(s)

	var s3toEncryption []S3ToEncryption
	for _, bucket := range b {
		params := &s3.GetBucketEncryptionInput{
			Bucket: aws.String(*bucket.Name),
		}
		_, err := svc.GetBucketEncryption(context.TODO(), params)
		if err != nil && !strings.Contains(err.Error(), "ServerSideEncryptionConfigurationNotFoundError") {
			panic(err)
		} else if err != nil {
			s3toEncryption = append(s3toEncryption, S3ToEncryption{*bucket.Name, false})
		} else {
			s3toEncryption = append(s3toEncryption, S3ToEncryption{*bucket.Name, true})
		}
	}
	logger.Debug(fmt.Sprintf("%v", s3toEncryption))
	return s3toEncryption
}

type S3ToVersioning struct {
	BucketName string
	Versioning bool
}

func GetS3ToVersioning(s aws.Config, b []types.Bucket) []S3ToVersioning {
	logger.Debug("Getting list of S3 buckets not in region")
	svc := s3.NewFromConfig(s)

	var s3toVersioning []S3ToVersioning
	for _, bucket := range b {
		params := &s3.GetBucketVersioningInput{
			Bucket: aws.String(*bucket.Name),
		}
		resp, err := svc.GetBucketVersioning(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if resp.Status != types.BucketVersioningStatusEnabled {
			s3toVersioning = append(s3toVersioning, S3ToVersioning{*bucket.Name, false})
		} else {
			s3toVersioning = append(s3toVersioning, S3ToVersioning{*bucket.Name, true})
		}
	}
	logger.Debug(fmt.Sprintf("%v", s3toVersioning))
	return s3toVersioning
}
