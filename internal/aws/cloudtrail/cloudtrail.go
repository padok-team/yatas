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
	check.Name = "Cloudtrails Encryption"
	check.Id = testName
	check.Description = "Check if all cloudtrails are encrypted"
	check.Status = "OK"
	for _, cloudtrail := range cloudtrails {
		if cloudtrail.KmsKeyId == nil || *cloudtrail.KmsKeyId == "" {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Cloudtrail " + *cloudtrail.Name + " is not encrypted"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *cloudtrail.TrailARN})
		} else {
			status := "OK"
			Message := "Cloudtrail " + *cloudtrail.Name + " is encrypted"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *cloudtrail.TrailARN})
		}
	}
	queueToAdd <- check
}

func CheckIfCloudtrailsGlobalServiceEventsEnabled(wg *sync.WaitGroup, s aws.Config, cloudtrails []types.Trail, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "Cloudtrails Global Service Events Activated"
	check.Id = testName
	check.Description = "Check if all cloudtrails have global service events enabled"
	check.Status = "OK"
	for _, cloudtrail := range cloudtrails {
		if !*cloudtrail.IncludeGlobalServiceEvents {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Cloudtrail " + *cloudtrail.Name + " has global service events disabled"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *cloudtrail.TrailARN})
		} else {
			status := "OK"
			Message := "Cloudtrail " + *cloudtrail.Name + " has global service events enabled"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *cloudtrail.TrailARN})
		}
	}
	queueToAdd <- check
}

func CheckIfCloudtrailsMultiRegion(wg *sync.WaitGroup, s aws.Config, cloudtrails []types.Trail, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "Cloudtrails Multi Region"
	check.Id = testName
	check.Description = "Check if all cloudtrails are multi region"
	check.Status = "OK"
	for _, cloudtrail := range cloudtrails {
		if !*cloudtrail.IsMultiRegionTrail {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Cloudtrail " + *cloudtrail.Name + " is not multi region"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *cloudtrail.TrailARN})
		} else {
			status := "OK"
			Message := "Cloudtrail " + *cloudtrail.Name + " is multi region"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *cloudtrail.TrailARN})
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
