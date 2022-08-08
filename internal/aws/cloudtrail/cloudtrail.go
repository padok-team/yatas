package cloudtrail

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func GetCloudtrails(s *session.Session) []*cloudtrail.Trail {
	svc := cloudtrail.New(s)
	input := &cloudtrail.DescribeTrailsInput{
		IncludeShadowTrails: aws.Bool(true),
	}
	result, err := svc.DescribeTrails(input)
	if err != nil {
		panic(err)
	}
	return result.TrailList
}

func CheckIfCloudtrailsEncrypted(s *session.Session, cloudtrails []*cloudtrail.Trail, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))

	var check types.Check
	check.Name = "Cloudtrails Encryption"
	check.Id = testName
	check.Description = "Check if all cloudtrails are encrypted"
	check.Status = "OK"
	for _, cloudtrail := range cloudtrails {
		if cloudtrail.KmsKeyId == nil || *cloudtrail.KmsKeyId == "" {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Cloudtrail " + *cloudtrail.Name + " is not encrypted"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *cloudtrail.TrailARN})
		} else {
			status := "OK"
			Message := "Cloudtrail " + *cloudtrail.Name + " is encrypted"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *cloudtrail.TrailARN})
		}
	}
	*c = append(*c, check)
}

func CheckIfCloudtrailsGlobalServiceEventsEnabled(s *session.Session, cloudtrails []*cloudtrail.Trail, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "Cloudtrails Global Service Events Activated"
	check.Id = testName
	check.Description = "Check if all cloudtrails have global service events enabled"
	check.Status = "OK"
	for _, cloudtrail := range cloudtrails {
		if !*cloudtrail.IncludeGlobalServiceEvents {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Cloudtrail " + *cloudtrail.Name + " has global service events disabled"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *cloudtrail.TrailARN})
		} else {
			status := "OK"
			Message := "Cloudtrail " + *cloudtrail.Name + " has global service events enabled"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *cloudtrail.TrailARN})
		}
	}
	*c = append(*c, check)
}

func CheckIfCloudtrailsMultiRegion(s *session.Session, cloudtrails []*cloudtrail.Trail, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "Cloudtrails Multi Region"
	check.Id = testName
	check.Description = "Check if all cloudtrails are multi region"
	check.Status = "OK"
	for _, cloudtrail := range cloudtrails {
		if !*cloudtrail.IsMultiRegionTrail {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Cloudtrail " + *cloudtrail.Name + " is not multi region"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *cloudtrail.TrailARN})
		} else {
			status := "OK"
			Message := "Cloudtrail " + *cloudtrail.Name + " is multi region"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *cloudtrail.TrailARN})
		}
	}
	*c = append(*c, check)
}

func RunCloudtrailTests(s *session.Session, c *config.Config) []types.Check {
	var checks []types.Check
	cloudtrails := GetCloudtrails(s)
	config.CheckTest(c, "AWS_CLD_001", CheckIfCloudtrailsEncrypted)(s, cloudtrails, "AWS_CLD_001", &checks)
	config.CheckTest(c, "AWS_CLD_002", CheckIfCloudtrailsGlobalServiceEventsEnabled)(s, cloudtrails, "AWS_CLD_002", &checks)
	config.CheckTest(c, "AWS_CLD_003", CheckIfCloudtrailsMultiRegion)(s, cloudtrails, "AWS_CLD_003", &checks)
	return checks
}
