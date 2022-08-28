package rds

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/stangirard/yatas/internal/yatas"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []yatas.Check) {

	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []yatas.Check
	svc := rds.NewFromConfig(s)

	instances := GetListRDS(svc)

	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_001", checkIfEncryptionEnabled)(checkConfig, instances, "AWS_RDS_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_002", checkIfBackupEnabled)(checkConfig, instances, "AWS_RDS_002")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_003", checkIfAutoUpgradeEnabled)(checkConfig, instances, "AWS_RDS_003")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_004", checkIfRDSPrivateEnabled)(checkConfig, instances, "AWS_RDS_004")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_005", CheckIfLoggingEnabled)(checkConfig, instances, "AWS_RDS_005")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_006", CheckIfDeleteProtectionEnabled)(checkConfig, instances, "AWS_RDS_006")
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
