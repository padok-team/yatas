package cloudfront

import (
	"fmt"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfStandardLogginEnabled(checkConfig yatas.CheckConfig, d []SummaryToConfig, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Cloudfronts queries are logged", "Check if all cloudfront distributions have standard logging enabled", testName)
	for _, cc := range d {

		if cc.config.Logging != nil && cc.config.Logging.Enabled != nil && *cc.config.Logging.Enabled {
			Message := "Standard logging is enabled on " + *cc.summary.Id
			result := results.Result{Status: "OK", Message: Message, ResourceID: *cc.summary.Id}
			check.AddResult(result)
		} else {
			Message := "Standard logging is not enabled on " + *cc.summary.Id
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *cc.summary.Id}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
