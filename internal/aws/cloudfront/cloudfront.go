package cloudfront

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfCookieLogginEnabled(checkConfig yatas.CheckConfig, d []SummaryToConfig, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Cookies Logging Enabled", "Check if all cloudfront distributions have cookies logging enabled", testName)
	for _, cc := range d {
		if cc.config.Logging != nil && *cc.config.Logging.Enabled && *cc.config.Logging.IncludeCookies {
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

func CheckIfACLUsed(checkConfig yatas.CheckConfig, d []SummaryToConfig, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("ACL Used", "Check if all cloudfront distributions have an ACL used", testName)
	for _, cc := range d {

		if *cc.config.WebACLId != "" {
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

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []results.Check
	d := GetAllCloudfront(s)
	s2c := GetAllDistributionConfig(s, d)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_CFT_001", CheckIfCloudfrontTLS1_2Minimum)(checkConfig, d, "AWS_CFT_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_CFT_002", CheckIfHTTPSOnly)(checkConfig, d, "AWS_CFT_002")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_CFT_003", CheckIfStandardLogginEnabled)(checkConfig, s2c, "AWS_CFT_003")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_CFT_004", CheckIfCookieLogginEnabled)(checkConfig, s2c, "AWS_CFT_004")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_CFT_005", CheckIfACLUsed)(checkConfig, s2c, "AWS_CFT_005")

	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
