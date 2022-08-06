package cloudtrail

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
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

func CheckIfCloudtrailsEncrypted(s *session.Session, cloudtrails []*cloudtrail.Trail, c *[]types.Check) {
	logger.Info("Running AWS_CLD_001")
	var check types.Check
	check.Name = "Cloudtrails Encryption"
	check.Id = "AWS_CLD_001"
	check.Description = "Check if all cloudtrails are encrypted"
	check.Status = "OK"
	for _, cloudtrail := range cloudtrails {
		if *cloudtrail.KmsKeyId != "" {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Cloudtrail " + *cloudtrail.Name + " is not encrypted"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "Cloudtrail " + *cloudtrail.Name + " is encrypted"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func CheckIfCloudtrailsGlobalServiceEventsEnabled(s *session.Session, cloudtrails []*cloudtrail.Trail, c *[]types.Check) {
	logger.Info("Running AWS_CLD_002")
	var check types.Check
	check.Name = "Cloudtrails Global Service Events Activated"
	check.Id = "AWS_CLD_002"
	check.Description = "Check if all cloudtrails have global service events enabled"
	check.Status = "OK"
	for _, cloudtrail := range cloudtrails {
		if *cloudtrail.IncludeGlobalServiceEvents != true {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Cloudtrail " + *cloudtrail.Name + " has global service events disabled"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "Cloudtrail " + *cloudtrail.Name + " has global service events enabled"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func CheckIfCloudtrailsMultiRegion(s *session.Session, cloudtrails []*cloudtrail.Trail, c *[]types.Check) {
	logger.Info("Running AWS_CLD_003")
	var check types.Check
	check.Name = "Cloudtrails Multi Region"
	check.Id = "AWS_CLD_003"
	check.Description = "Check if all cloudtrails are multi region"
	check.Status = "OK"
	for _, cloudtrail := range cloudtrails {
		if *cloudtrail.IsMultiRegionTrail != true {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Cloudtrail " + *cloudtrail.Name + " is not multi region"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "Cloudtrail " + *cloudtrail.Name + " is multi region"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func RunCloudtrailTests(s *session.Session) []types.Check {
	var checks []types.Check
	cloudtrails := GetCloudtrails(s)
	CheckIfCloudtrailsEncrypted(s, cloudtrails, &checks)
	CheckIfCloudtrailsGlobalServiceEventsEnabled(s, cloudtrails, &checks)
	CheckIfCloudtrailsMultiRegion(s, cloudtrails, &checks)
	return checks
}
