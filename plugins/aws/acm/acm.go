package acm

import (
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/stangirard/yatas/internal/yatas"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []yatas.Check) {
	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []yatas.Check
	svc := acm.NewFromConfig(s)
	certificates := GetCertificates(svc)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_ACM_001", CheckIfACMValid)(checkConfig, certificates, "AWS_ACM_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_ACM_002", CheckIfCertificateExpiresIn90Days)(checkConfig, certificates, "AWS_ACM_002")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_ACM_003", CheckIfACMInUse)(checkConfig, certificates, "AWS_ACM_003")
	go func() {
		for t := range checkConfig.Queue {
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
