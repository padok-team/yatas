package ecr

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
	ecr := GetECRs(s)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_ECR_001", CheckIfImageScanningEnabled)(checkConfig, ecr, "AWS_ECR_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_ECR_002", CheckIfEncrypted)(checkConfig, ecr, "AWS_ECR_002")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_ECR_003", CheckIfTagImmutable)(checkConfig, ecr, "AWS_ECR_003")
	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			if c.ProgressDetailed != nil {
				c.ProgressDetailed.Increment()
				time.Sleep(time.Millisecond * 100)
			}
			checkConfig.Wg.Done()

		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
