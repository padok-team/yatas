package loadbalancers

import (
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/internal/yatas"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []yatas.Check) {

	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []yatas.Check
	loadBalancers := GetElasticLoadBalancers(s)
	la := GetLoadBalancersAttributes(s, loadBalancers)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_LB_001", CheckIfAccessLogsEnabled)(checkConfig, la, "AWS_ELB_001")
	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)
			if c.CheckProgress.Bar != nil {
				c.CheckProgress.Bar.Increment()
				time.Sleep(time.Millisecond * 100)
			}

			checkConfig.Wg.Done()

		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
