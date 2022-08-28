package ecr

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/internal/yatas"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []yatas.Check) {

	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []yatas.Check
	ecr := GetECRs(s)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_ECR_001", CheckIfImageScanningEnabled)(checkConfig, ecr, "AWS_ECR_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_ECR_002", CheckIfEncrypted)(checkConfig, ecr, "AWS_ECR_002")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_ECR_003", CheckIfTagImmutable)(checkConfig, ecr, "AWS_ECR_003")
	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)
			if c.CheckProgress.Bar != nil {
				c.CheckProgress.Bar.Increment()
			}

			checkConfig.Wg.Done()

		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
