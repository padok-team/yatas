package guardduty

import (
	"fmt"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfGuarddutyEnabled(checkConfig yatas.CheckConfig, testName string, detectors []string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("GuardDuty is enabled in the account", "Check if GuardDuty is enabled", testName)

	if len(detectors) == 0 {
		Message := "GuardDuty is not enabled"
		result := yatas.Result{Status: "FAIL", Message: Message}
		check.AddResult(result)
	} else {
		Message := "GuardDuty is enabled"
		result := yatas.Result{Status: "OK", Message: Message}
		check.AddResult(result)
	}
	checkConfig.Queue <- check
}
