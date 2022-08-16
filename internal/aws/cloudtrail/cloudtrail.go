package cloudtrail

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func GetCloudtrails(s aws.Config) []types.Trail {
	svc := cloudtrail.NewFromConfig(s)
	input := &cloudtrail.DescribeTrailsInput{
		IncludeShadowTrails: aws.Bool(true),
	}
	result, err := svc.DescribeTrails(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	return result.TrailList
}

func CheckIfCloudtrailsEncrypted(checkConfig yatas.CheckConfig, cloudtrails []types.Trail, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))

	var check results.Check
	check.InitCheck("Cloudtrails Encryption", "check if all cloudtrails are encrypted", testName)
	for _, cloudtrail := range cloudtrails {
		if cloudtrail.KmsKeyId == nil || *cloudtrail.KmsKeyId == "" {
			Message := "Cloudtrail " + *cloudtrail.Name + " is not encrypted"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		} else {
			Message := "Cloudtrail " + *cloudtrail.Name + " is encrypted"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

func CheckIfCloudtrailsGlobalServiceEventsEnabled(checkConfig yatas.CheckConfig, cloudtrails []types.Trail, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Cloudtrails Global Service Events Activated", "check if all cloudtrails have global service events enabled", testName)
	for _, cloudtrail := range cloudtrails {
		if !*cloudtrail.IncludeGlobalServiceEvents {
			Message := "Cloudtrail " + *cloudtrail.Name + " has global service events disabled"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		} else {
			Message := "Cloudtrail " + *cloudtrail.Name + " has global service events enabled"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

func CheckIfCloudtrailsMultiRegion(checkConfig yatas.CheckConfig, cloudtrails []types.Trail, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Cloudtrails Multi Region", "check if all cloudtrails are multi region", testName)
	for _, cloudtrail := range cloudtrails {
		if !*cloudtrail.IsMultiRegionTrail {
			Message := "Cloudtrail " + *cloudtrail.Name + " is not multi region"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		} else {
			Message := "Cloudtrail " + *cloudtrail.Name + " is multi region"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {
	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []results.Check
	cloudtrails := GetCloudtrails(s)

	go yatas.CheckTest(checkConfig.Wg, c, "AWS_CLD_001", CheckIfCloudtrailsEncrypted)(checkConfig, cloudtrails, "AWS_CLD_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_CLD_002", CheckIfCloudtrailsGlobalServiceEventsEnabled)(checkConfig, cloudtrails, "AWS_CLD_002")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_CLD_003", CheckIfCloudtrailsMultiRegion)(checkConfig, cloudtrails, "AWS_CLD_003")

	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks

}
