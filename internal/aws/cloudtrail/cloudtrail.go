package cloudtrail

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

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
