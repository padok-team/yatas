package vpc

import (
	"fmt"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func checkIfVPCFLowLogsEnabled(checkConfig yatas.CheckConfig, VpcFlowLogs []VpcToFlowLogs, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("VPC Flow Logs are activated", "Check if VPC Flow Logs are enabled", testName)
	for _, vpcFlowLog := range VpcFlowLogs {

		if len(vpcFlowLog.FlowLogs) == 0 {
			Message := "VPC Flow Logs are not enabled on " + vpcFlowLog.VpcID
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: vpcFlowLog.VpcID}
			check.AddResult(result)
		} else {
			Message := "VPC Flow Logs are enabled on " + vpcFlowLog.VpcID
			result := yatas.Result{Status: "OK", Message: Message, ResourceID: vpcFlowLog.VpcID}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
