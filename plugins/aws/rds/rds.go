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
	clusters := GetListDBClusters(svc)

	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_001", checkIfEncryptionEnabled)(checkConfig, instances, "AWS_RDS_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_002", checkIfBackupEnabled)(checkConfig, instances, "AWS_RDS_002")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_003", checkIfAutoUpgradeEnabled)(checkConfig, instances, "AWS_RDS_003")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_004", checkIfRDSPrivateEnabled)(checkConfig, instances, "AWS_RDS_004")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_005", CheckIfLoggingEnabled)(checkConfig, instances, "AWS_RDS_005")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_006", CheckIfDeleteProtectionEnabled)(checkConfig, instances, "AWS_RDS_006")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_007", checkIfClusterAutoUpgradeEnabled)(checkConfig, clusters, "AWS_RDS_007")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_008", checkIfClusterBackupEnabled)(checkConfig, clusters, "AWS_RDS_008")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_009", CheckIfClusterDeleteProtectionEnabled)(checkConfig, clusters, "AWS_RDS_009")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_010", checkIfClusterEncryptionEnabled)(checkConfig, clusters, "AWS_RDS_010")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_011", CheckIfClusterLoggingEnabled)(checkConfig, clusters, "AWS_RDS_011")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_012", checkIfClusterRDSPrivateEnabled)(checkConfig, clusters, "AWS_RDS_012")

	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)

			checkConfig.Wg.Done()

		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
