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

func CheckIfCloudtrailsEncrypted(wg *sync.WaitGroup, s aws.Config, cloudtrails []types.Trail, testName string, queueToAdd chan results.Check) {
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
	queueToAdd <- check
}

func CheckIfCloudtrailsGlobalServiceEventsEnabled(wg *sync.WaitGroup, s aws.Config, cloudtrails []types.Trail, testName string, queueToAdd chan results.Check) {
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
	queueToAdd <- check
}

func CheckIfCloudtrailsMultiRegion(wg *sync.WaitGroup, s aws.Config, cloudtrails []types.Trail, testName string, queueToAdd chan results.Check) {
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
	queueToAdd <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {
	var checks []results.Check
	var wg sync.WaitGroup
	queueResults := make(chan results.Check, 10)
	cloudtrails := GetCloudtrails(s)

	go yatas.CheckTest(&wg, c, "AWS_CLD_001", CheckIfCloudtrailsEncrypted)(&wg, s, cloudtrails, "AWS_CLD_001", queueResults)
	go yatas.CheckTest(&wg, c, "AWS_CLD_002", CheckIfCloudtrailsGlobalServiceEventsEnabled)(&wg, s, cloudtrails, "AWS_CLD_002", queueResults)
	go yatas.CheckTest(&wg, c, "AWS_CLD_003", CheckIfCloudtrailsMultiRegion)(&wg, s, cloudtrails, "AWS_CLD_003", queueResults)

	go func() {
		for t := range queueResults {
			checks = append(checks, t)
			wg.Done()
		}
	}()

	wg.Wait()

	queue <- checks

}
