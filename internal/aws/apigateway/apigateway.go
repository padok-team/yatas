package apigateway

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
	apis := GetApiGateways(s)
	stages := GetAllStagesApiGateway(s, apis)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_APG_001", CheckIfStagesCloudwatchLogsExist)(checkConfig, stages, "AWS_APG_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_APG_002", CheckIfStagesProtectedByAcl)(checkConfig, stages, "AWS_APG_002")

	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
