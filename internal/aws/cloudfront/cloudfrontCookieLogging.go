package cloudfront

import (
	"fmt"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfCookieLogginEnabled(checkConfig yatas.CheckConfig, d []SummaryToConfig, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Cookies Logging Enabled", "Check if all cloudfront distributions have cookies logging enabled", testName)
	for _, cc := range d {
		if cc.config.Logging != nil && *cc.config.Logging.Enabled && cc.config.Logging.IncludeCookies != nil && *cc.config.Logging.IncludeCookies {
			Message := "Cookie logging is enabled on " + *cc.summary.Id
			result := results.Result{Status: "OK", Message: Message, ResourceID: *cc.summary.Id}
			check.AddResult(result)
		} else {
			Message := "Cookie logging is not enabled on " + *cc.summary.Id
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *cc.summary.Id}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
