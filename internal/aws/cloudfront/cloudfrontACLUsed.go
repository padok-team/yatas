package cloudfront

import (
	"fmt"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfACLUsed(checkConfig yatas.CheckConfig, d []SummaryToConfig, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("ACL Used", "Check if all cloudfront distributions have an ACL used", testName)
	for _, cc := range d {

		if cc.config.WebACLId != nil && *cc.config.WebACLId != "" {
			Message := "ACL is used on " + *cc.summary.Id
			result := results.Result{Status: "OK", Message: Message, ResourceID: *cc.summary.Id}
			check.AddResult(result)
		} else {
			Message := "ACL is not used on " + *cc.summary.Id
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *cc.summary.Id}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
