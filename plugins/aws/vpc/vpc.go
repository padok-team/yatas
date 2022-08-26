package vpc

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/internal/yatas"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []yatas.Check) {

	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []yatas.Check
	vpcs := GetListVPC(s)
	subnetsforvpcs := GetSubnetForVPCS(s, vpcs)
	internetGatewaysForVpc := GetInternetGatewaysForVpc(s, vpcs)
	vpcFlowLogs := GetFlowLogsForVpc(s, vpcs)

	go yatas.CheckTest(checkConfig.Wg, c, "AWS_VPC_001", checkCIDR20)(checkConfig, vpcs, "AWS_VPC_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_VPC_002", checkIfOnlyOneVPC)(checkConfig, vpcs, "AWS_VPC_002")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_VPC_003", checkIfOnlyOneGateway)(checkConfig, internetGatewaysForVpc, "AWS_VPC_003")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_VPC_004", checkIfVPCFLowLogsEnabled)(checkConfig, vpcFlowLogs, "AWS_VPC_004")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_VPC_005", CheckIfAtLeast2Subnets)(checkConfig, subnetsforvpcs, "AWS_VPC_005")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_VPC_006", CheckIfSubnetInDifferentZone)(checkConfig, subnetsforvpcs, "AWS_VPC_006")
	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
