package cloudfront

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfHTTPSOnly(checkConfig yatas.CheckConfig, d []types.DistributionSummary, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Cloudfronts only allow HTTPS or redirect to HTTPS", "Check if all cloudfront distributions are HTTPS only", testName)
	for _, cloudfront := range d {
		if cloudfront.DefaultCacheBehavior != nil && (cloudfront.DefaultCacheBehavior.ViewerProtocolPolicy == "https-only" || cloudfront.DefaultCacheBehavior.ViewerProtocolPolicy == "redirect-to-https") {
			Message := "Cloudfront distribution is HTTPS only on " + *cloudfront.Id
			result := results.Result{Status: "OK", Message: Message, ResourceID: *cloudfront.Id}
			check.AddResult(result)
		} else {
			Message := "Cloudfront distribution is not HTTPS only on " + *cloudfront.Id
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *cloudfront.Id}
			check.AddResult(result)
		}
	}

	checkConfig.Queue <- check
}
