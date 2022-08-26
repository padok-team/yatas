package dynamodb

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
	dynamodbs := GetDynamodbs(s)
	gt := GetTables(s, dynamodbs)
	gb := GetContinuousBackups(s, dynamodbs)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_DYN_001", CheckIfDynamodbEncrypted)(checkConfig, gt, "AWS_DYN_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_DYN_002", CheckIfDynamodbContinuousBackupsEnabled)(checkConfig, gb, "AWS_DYN_002")

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
