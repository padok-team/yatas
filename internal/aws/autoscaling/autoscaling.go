package autoscaling

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {
	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []results.Check
	svc := autoscaling.NewFromConfig(s)
	groups := GetAutoscalingGroups(svc)

	go yatas.CheckTest(checkConfig.Wg, c, "AWS_ASG_001", CheckIfDesiredCapacityMaxCapacityBelow80percent)(checkConfig, groups, "AWS_ASG_001")

	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
