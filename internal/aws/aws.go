package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/internal/aws/apigateway"
	"github.com/stangirard/yatas/internal/aws/autoscaling"
	"github.com/stangirard/yatas/internal/aws/cloudfront"
	"github.com/stangirard/yatas/internal/aws/cloudtrail"
	"github.com/stangirard/yatas/internal/aws/dynamodb"
	"github.com/stangirard/yatas/internal/aws/ec2"
	"github.com/stangirard/yatas/internal/aws/ecr"
	"github.com/stangirard/yatas/internal/aws/iam"
	"github.com/stangirard/yatas/internal/aws/lambda"
	"github.com/stangirard/yatas/internal/aws/loadbalancers"
	"github.com/stangirard/yatas/internal/aws/rds"
	"github.com/stangirard/yatas/internal/aws/s3"
	"github.com/stangirard/yatas/internal/aws/volumes"
	"github.com/stangirard/yatas/internal/aws/vpc"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func Run(c *yatas.Config) ([]results.Check, error) {
	s := initAuth(c)
	logger.Info("Launching AWS checks")
	checks := initTest(s, c)
	return checks, nil
}

func initTest(s aws.Config, c *yatas.Config) []results.Check {

	var checks []results.Check
	checks = append(checks, s3.RunChecks(s, c)...)
	checks = append(checks, volumes.RunChecks(s, c)...)
	checks = append(checks, rds.RunChecks(s, c)...)
	checks = append(checks, vpc.RunChecks(s, c)...)
	checks = append(checks, cloudtrail.RunChecks(s, c)...)
	checks = append(checks, ecr.RunChecks(s, c)...)
	checks = append(checks, lambda.RunChecks(s, c)...)
	checks = append(checks, dynamodb.RunChecks(s, c)...)
	checks = append(checks, ec2.RunChecks(s, c)...)
	checks = append(checks, iam.RunChecks(s, c)...)
	checks = append(checks, cloudfront.RunChecks(s, c)...)
	checks = append(checks, apigateway.RunChecks(s, c)...)
	checks = append(checks, autoscaling.RunChecks(s, c)...)
	checks = append(checks, loadbalancers.RunChecks(s, c)...)
	logger.Info("AWS checks completed ✅")

	return checks
}
