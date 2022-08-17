package guardduty

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func GetDetectors(s aws.Config) []string {
	svc := guardduty.NewFromConfig(s)
	input := &guardduty.ListDetectorsInput{}
	result, err := svc.ListDetectors(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	return result.DetectorIds
}

func CheckIfGuarddutyEnabled(checkConfig yatas.CheckConfig, testName string, detectors []string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("GuardDuty Enabled", "Check if GuardDuty is enabled", testName)

	if len(detectors) == 0 {
		Message := "GuardDuty is not enabled"
		result := results.Result{Status: "FAIL", Message: Message}
		check.AddResult(result)
	} else {
		Message := "GuardDuty is enabled"
		result := results.Result{Status: "OK", Message: Message}
		check.AddResult(result)
	}
	checkConfig.Queue <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []results.Check
	guardyDetectors := GetDetectors(checkConfig.ConfigAWS)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_GDT_001", CheckIfGuarddutyEnabled)(checkConfig, "AWS_GDT_001", guardyDetectors)
	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
