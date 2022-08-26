package vpc

import (
	"fmt"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfAtLeast2Subnets(checkConfig yatas.CheckConfig, vpcToSubnets []VPCToSubnet, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("VPC have at least 2 subnets", "Check if VPC has at least 2 subnets", testName)
	for _, vpcToSubnet := range vpcToSubnets {

		if len(vpcToSubnet.Subnets) < 2 {
			Message := "VPC " + vpcToSubnet.VpcID + " has less than 2 subnets"
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: vpcToSubnet.VpcID}
			check.AddResult(result)
		} else {
			Message := "VPC " + vpcToSubnet.VpcID + " has at least 2 subnets"
			result := yatas.Result{Status: "OK", Message: Message, ResourceID: vpcToSubnet.VpcID}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
